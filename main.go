package main

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/Ikenna-Okpala/personal-website-revamp.git/internal/view"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)


func Mailer() * smtp.Auth{

	user:= os.Getenv("EMAIL_USER")
	password:= os.Getenv("EMAIL_PASSWORD")

	host:= os.Getenv("EMAIL_HOST")



	auth:= smtp.PlainAuth("", user, password, host)

	return &auth
}

func Db() * datastore.Client{

	ctx:= context.Background()

	client, err:= datastore.NewClient(ctx, "giving-talent")

	if err != nil {
		log.Fatalln(err)
	}

	return client
}

func Blob() * storage.Client {

	ctx:= context.Background()

	client, err:= storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	return client
}

type Server struct{
	Mailer * smtp.Auth
	Db * datastore.Client
	Blob * storage.Client
}



func (s Server) Email(w http.ResponseWriter, r * http.Request){

	
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	msg:= r.FormValue("message")

	name:= r.FormValue("name")

	email:= r.FormValue("email")

	if msg == "" || name == "" || email == ""{
		http.Error(w, "Form fields invalid", http.StatusBadRequest)
		return
	}

	html:= bytes.NewBufferString("")

	view.Email(view.Contact{Name: name, Email: email, Message: msg}).Render(r.Context(), html)

	

	host:= os.Getenv("EMAIL_HOST")

	port:= os.Getenv("EMAIL_PORT")

	from:= os.Getenv("EMAIL_USER")

	to:= []string{os.Getenv("MY_EMAIL")}

	subject:= "Subject: SOMEONE HAS REACHED OUT TO YOU\n"
	mime:= "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\"; \n\n"
	body:= html.String()

	msgByte:= []byte(subject + mime + body)


	err2:= smtp.SendMail(host+":"+port, *s.Mailer, from, to, msgByte)

	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}




	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err:= io.WriteString(w, "success"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {return true},
}

func (s Server) LiveReload(w http.ResponseWriter, r * http.Request){
	conn, err:= upgrader.Upgrade(w, r, nil)

	if err != nil{
		panic(err)
	}

	defer conn.Close()

	for {
		conn.ReadMessage()
	}
}



func (s Server) AddProject(w http.ResponseWriter, r * http.Request){

	err:= r.ParseMultipartForm(32 << 20)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, _, err3 := r.FormFile("screenshot")

	if err3 != nil{

		http.Error(w, err3.Error(), http.StatusBadRequest)
		return
	}

	req := r.FormValue("json")
	if req == ""{
		http.Error(w, "Request not provided", http.StatusBadRequest)
		return
	}
	

	var p view.Project

	err = json.Unmarshal([] byte(req), &p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bucketName:= os.Getenv("GCP_BUCKET_NAME")

	if bucketName == ""{
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}


	objName:= fmt.Sprintf("%s.png", p.Name)

	obj:= s.Blob.Bucket(bucketName).Object(objName)

	writer := obj.NewWriter(r.Context())

	if _, err:= io.Copy(writer, img); err != nil{
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if err = writer.Close(); err!= nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	p.ImageUrl = fmt.Sprintf("https://storage.googleapis.com/personal-website-golang/%s", objName)

	
	// key:= datastore.IncompleteKey("Project", nil)
	id:= uuid.New().String()

	key:= datastore.NameKey("Project", id, nil)

	_, err2:= s.Db.Put(r.Context(), key, &p)

	if err2 != nil{
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	p.ID = key.Name

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(&p)
}

func (s Server) DeleteProject(w http.ResponseWriter, r * http.Request){

	id:= r.PathValue("id")

	if id == ""{
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}



	bucketName:= os.Getenv("GCP_BUCKET_NAME")

	if bucketName == ""{
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	key := datastore.NameKey("Project", id, nil)

	p:= new(view.Project)

	if err:= s.Db.Get(r.Context(), key, p); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return 
	}

	index:= strings.LastIndex(p.ImageUrl, "/")

	if index == -1{
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// provided that the file name does not include a slash

	objectName:= p.ImageUrl[index + 1:]

	if err:= s.Blob.Bucket(bucketName).Object(objectName).Delete(r.Context()); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}


	if err := s.Db.Delete(r.Context(), key); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s Server) UpdateProject(w http.ResponseWriter, r * http.Request){

	id:= r.PathValue("id")

	if id == ""{
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	err:= r.ParseMultipartForm(32 << 20)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img, _, err3 := r.FormFile("screenshot")

	if err3 != nil{

		http.Error(w, err3.Error(), http.StatusBadRequest)
		return
	}

	req := r.FormValue("json")
	if req == ""{
		http.Error(w, "Request not provided", http.StatusBadRequest)
		return
	}

	newP:= new(view.Project)

	if err = json.Unmarshal([] byte(req), &newP); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bucketName:= os.Getenv("GCP_BUCKET_NAME")

	objName:= fmt.Sprintf("%s.png", newP.Name)

	obj:= s.Blob.Bucket(bucketName).Object(objName)

	writer := obj.NewWriter(r.Context())

	if _, err:= io.Copy(writer, img); err != nil{
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if err = writer.Close(); err!= nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	key:= datastore.NameKey("Project", id , nil)

	p:= new(view.Project)

	newP.ImageUrl = fmt.Sprintf("https://storage.googleapis.com/personal-website-golang/%s", objName)


	if err:= s.Db.Get(r.Context(), key, p); err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	*p = *newP

	if _, err:= s.Db.Put(r.Context(), key, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.ID = key.Name

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	
	json.NewEncoder(w).Encode(p)
}

func (s Server ) GetProjects() ([]view.Project, error) {

	var projects []*view.Project

	keys, err := s.Db.GetAll(context.Background(), datastore.NewQuery("Project").Order("-FinishedAt"), &projects)

	if err != nil {
		return nil, err
	}

	var result []view.Project

	for i, p := range projects {

		p.ID = keys[i].Name
		result = append(result, *p)
	}

	return result, nil

}


func (s Server) GetProjectsApi(w http.ResponseWriter, r * http.Request){

	projects, err:= s.GetProjects()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&projects)

}

func (s Server) Home(w http.ResponseWriter, r * http.Request){

	projects, err:= s.GetProjects()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hx:= r.Header.Get("HX-Request")

	if hx == ""{
		view.App("Ikenna's Portfolio", view.About(projects)).Render(r.Context(), w)
	}else{
		view.About(projects).Render(r.Context(), w)
	}

	
}


func (s Server) Projects(w http.ResponseWriter, r* http.Request){

	category:= r.URL.Query().Get("category")
	skill := r.URL.Query().Get("skill")

	if category != ""{
		if category == "all"{

			projects, err:= s.GetProjects()
	
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	
			view.ProjectCards(projects).Render(r.Context(), w)
			return
		} else {
	
			query:= datastore.NewQuery("Project").
			FilterField("Category", "=", category).
			Order("-FinishedAt")
	
			it:= s.Db.Run(r.Context(), query)
	
			projects:= [] view.Project{}
	
			for{
	
				p:= new(view.Project)
	
				_, err:= it.Next(p)
	
				if err == iterator.Done {
					break
				}
	
				if err!= nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
	
				projects = append(projects, *p)
			}
	
	
			view.ProjectCards(projects).Render(r.Context(), w)
			return
		}	
	}

	if skill != ""{
		query:= datastore.NewQuery("Project").
		FilterField("Skills", "=", skill).
		Order("-FinishedAt")

		it:= s.Db.Run(r.Context(), query)

		projects:= [] view.Project{}

		for {
			p:= new(view.Project)

			_, err:= it.Next(p)

			if err == iterator.Done{
				break
			}

			if err!= nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			projects = append(projects, *p)
		}

		view.ProjectCards(projects).Render(r.Context(), w)
		return
	}

	http.Error(w, "Bad query", http.StatusBadRequest)

	
}



func (s Server) AddBlog(w http.ResponseWriter, r * http.Request){

	blog:= new(view.Blog)

	err:= json.NewDecoder(r.Body).Decode(blog)

	if err != nil {
		http.Error(w, "Bad query", http.StatusBadRequest)
		return
	}

	blog.CreatedAt = time.Now().Format(time.RFC3339)
	blog.LastUpdated = blog.CreatedAt


	uuid:= uuid.New().String()

	key:= datastore.NameKey("Blog", uuid, nil)

	if _ , err:= s.Db.Put(r.Context(), key, blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blog.Id = key.Name

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(blog)
}

func (s Server) GetBlogs(filter string) ([] view.Blog, error){

	filter = strings.ToLower(filter)

	blogs:= new([] view.Blog)

	query:= datastore.NewQuery("Blog").Order("-CreatedAt")

	keys, err:= s.Db.GetAll(context.Background(), query, blogs)

	if err != nil {
		return nil, err
	}

	filterBlog := [] view.Blog {}

	for i, blog:= range *blogs {

		blog.Id = keys[i].Name

		title:= strings.ToLower(blog.Title)

		if strings.Contains(title, filter){
			
			filterBlog = append(filterBlog, blog)
			continue
		}

		for _, label:= range blog.Labels {

			label = strings.ToLower(label)

			if strings.Contains(label, filter){
				filterBlog = append(filterBlog, blog)
				break
			}
		}


	}

	return filterBlog, nil
}

func (s Server) GetBlogsApi(w http.ResponseWriter, r * http.Request){


	filter:= r.URL.Query().Get("search")

	blogs, err:= s.GetBlogs(filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&blogs)
}

func (s Server) UpdateBlog(w http.ResponseWriter, r * http.Request){

	id:= r.PathValue("id")

	reqBlog:= new(view.Blog)

	if err:= json.NewDecoder(r.Body).Decode(reqBlog); err != nil{
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}



	if id == ""{
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	key:= datastore.NameKey("Blog", id, nil)

	blog:= new(view.Blog)

	if err:= s.Db.Get(r.Context(), key, blog); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if reqBlog.Title != ""{
		blog.Title = reqBlog.Title
	}

	if len(reqBlog.Labels) > 0{
		blog.Labels = reqBlog.Labels
	}

	if reqBlog.BlogUrl != ""{
		blog.BlogUrl = reqBlog.BlogUrl
	}

	blog.LastUpdated = time.Now().Format(time.RFC3339)
	
	if _, err:= s.Db.Put(r.Context(), key, blog); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blog.Id = key.Name

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(blog)

}

func (s Server) DeleteBlog(w http.ResponseWriter, r * http.Request){

	id:= r.PathValue("id")

	if id == ""{
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	key:= datastore.NameKey("Blog", id, nil)

	if err:= s.Db.Delete(r.Context(), key); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)


}

func (s Server) Blog(w http.ResponseWriter, r * http.Request){

	blogs, err:= s.GetBlogs("")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hx:= r.Header.Get("HX-Request")

	if hx == ""{
		view.App("Ikenna's Portfolio", view.BlogUI(blogs)).Render(r.Context(), w)

	}else{
	view.BlogUI(blogs).Render(r.Context(), w)
	}

	
}

func (s Server) BlogSearch(w http.ResponseWriter, r * http.Request){

	search:= r.URL.Query().Get("search")

	blogs, err:= s.GetBlogs(search)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.BlogList(blogs).Render(r.Context(), w)
	
}

func (s Server) BlogTheSimpleMindOfAProgrammer(w http.ResponseWriter, r * http.Request){

	view.App("Ikenna's Portfolio", view.BlogSimpleMindProgrammer()).Render(r.Context(), w)
}

func (s Server) Gallery(w http.ResponseWriter, r * http.Request){

	photos, err:= s.GetPhotos()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	

	slices.SortFunc(*photos, func(p1 view.Photo, p2 view.Photo) int{
		return -cmp.Compare(p1.CreatedAt, p2.CreatedAt)

	})

	hx:= r.Header.Get("HX-Request")

	if hx == ""{
		view.App("Ikenna's Portfolio", view.Gallery(*photos)).Render(r.Context(), w)
	}else{
		view.Gallery(*photos).Render(r.Context(), w)
	}
	
}



func (s Server) AddPhoto (w http.ResponseWriter, r * http.Request){

	if err:= r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	photo, _, err:= r.FormFile("photo")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metadata:= r.FormValue("metadata")

	if metadata == ""{
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	metaDataToAdd:= new(view.Photo)

	json.Unmarshal([]byte(metadata), metaDataToAdd)


	bucketName:= os.Getenv("GCP_BUCKET_NAME")

	if bucketName == ""{
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	key:= datastore.NameKey("Gallery", uuid.NewString(), nil)

	photoName:= fmt.Sprintf("%s.png", key.Name)

	writer:= s.Blob.Bucket(bucketName).Object(photoName).NewWriter(r.Context())

	if _, err:= io.Copy(writer, photo); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if err:= writer.Close(); err!= nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	metaDataToAdd.Id = key.Name
	metaDataToAdd.Url = fmt.Sprintf("https://storage.googleapis.com/personal-website-golang/%s", photoName)

	if _, err:= s.Db.Put(r.Context(), key, metaDataToAdd); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(metaDataToAdd)

}

func (s Server) UpdatePhoto(w http.ResponseWriter, r * http.Request){
	
	id:= r.PathValue("id")

	if id == ""{
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	metadataUpdated:= new(view.Photo)

	if err:= json.NewDecoder(r.Body).Decode(metadataUpdated); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	key:= datastore.NameKey("Gallery", id, nil)

	currentMetadata:= new(view.Photo)

	if err:= s.Db.Get(r.Context(), key, currentMetadata); err!= nil {
		
		if errors.Is(err, datastore.ErrNoSuchEntity){
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if metadataUpdated.Caption != ""{
		currentMetadata.Caption = metadataUpdated.Caption
	}

	if metadataUpdated.CreatedAt != ""{
		currentMetadata.CreatedAt = metadataUpdated.CreatedAt
	}

	currentMetadata.Id = key.Name

	if _, err:= s.Db.Put(r.Context(), key, currentMetadata); err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(currentMetadata)
}

func (s Server) GetPhotos() ( * [] view.Photo, error) {

	photos := new([] view.Photo)

	query:= datastore.NewQuery("Gallery").
	Order("-CreatedAt")

	keys, err:= s.Db.GetAll(context.Background(), query, photos);

	if err != nil {
		return nil, errors.New("somthing went wrong")
	}

	for i, key:= range keys {

		(*photos)[i].Id = key.Name
	}

	return photos, nil

}

func (s Server) GetPhotosApi(w http.ResponseWriter, r * http.Request){

	photos, err:= s.GetPhotos()

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(photos)
}

func (s Server) DeletePhoto(w http.ResponseWriter, r * http.Request){

	id:= r.PathValue("id")

	if id == ""{
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	bucketName:= os.Getenv("GCP_BUCKET_NAME")

	if bucketName == ""{
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	key:= datastore.NameKey("Gallery", id, nil)

	objectName:= fmt.Sprintf("%s.png", id)
	
	if err:= s.Blob.Bucket(bucketName).Object(objectName).Delete(r.Context()); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}


	if err:= s.Db.Delete(r.Context(), key); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

}

func (s Server) Auth(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		username, password, ok:= r.BasicAuth()

		if ok {

			userHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))

			actualUserHash:= sha256.Sum256([]byte(os.Getenv("API_USERNAME")))
			actualPasswordHash:= sha256.Sum256([]byte(os.Getenv("API_PASSWORD")))

			isUserValid:= (subtle.ConstantTimeCompare(userHash[:], actualUserHash[:]) ) == 1

			isPasswordValid := (subtle.ConstantTimeCompare(passwordHash[:], actualPasswordHash[:])) == 1

			if isUserValid && isPasswordValid{
				next(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}


func main() {

	err:= godotenv.Load()

	if err != nil {
		log.Fatalln("Could not load ENV variables")
	}

	server:= Server{
		Mailer: Mailer(),
		Db: Db(),
		Blob: Blob(),
	}

	fs:= http.FileServer(http.Dir("./static"))


	http.HandleFunc("/", server.Home)

	http.HandleFunc("/live", server.LiveReload)

	http.HandleFunc("POST /api/v1/email", server.Email)
        

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("POST /api/v1/project", server.Auth(server.AddProject))

	http.HandleFunc("PUT /api/v1/project/{id}", server.Auth(server.UpdateProject))

	http.HandleFunc("DELETE /api/v1/project/{id}", server.Auth(server.DeleteProject))

	http.HandleFunc("GET /api/v1/project", server.Auth(server.GetProjectsApi))

	http.HandleFunc("GET /project", server.Projects)

	http.HandleFunc("POST /api/v1/blog", server.Auth(server.AddBlog))

	http.HandleFunc("GET /api/v1/blog", server.Auth(server.GetBlogsApi))

	http.HandleFunc("PATCH /api/v1/blog/{id}", server.Auth(server.UpdateBlog))

	http.HandleFunc("DELETE /api/v1/blog/{id}", server.Auth(server.DeleteBlog))

	http.HandleFunc("GET /blog", server.Blog)

	http.HandleFunc("GET /blog/search", server.BlogSearch)

	http.HandleFunc("GET /gallery", server.Gallery)

	http.HandleFunc("GET /blog/the-simple-mind-of-a-programmer", server.BlogTheSimpleMindOfAProgrammer)

	http.HandleFunc("GET /api/v1/photo", server.Auth(server.GetPhotosApi))

	http.HandleFunc("POST /api/v1/photo", server.Auth(server.AddPhoto))

	http.HandleFunc("PATCH /api/v1/photo/{id}", server.Auth(server.UpdatePhoto))

	http.HandleFunc("DELETE /api/v1/photo/{id}", server.Auth(server.DeletePhoto))

	fmt.Println("Listening on port: 8080")

	err2:= http.ListenAndServe("localhost:8080", nil)

	if err2 != nil {
		panic(err2)  
	}

	


	

}