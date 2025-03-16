// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819

package view

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

type Project struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Skills       []string `json:"skills"`
	ImageUrl     string   `json:"imageUrl"`
	GithubLink   string   `json:"githubLink"`
	Category     string   `json:"category"`
	ID           string   `json:"id"`
	ImageLinkUrl string   `json:"imageLinkUrl"`
	FinishedAt   string   `json:"finishedAt"`
}

func Footer() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<footer class=\"w-full flex flex-row justify-between px-64 h-32\"><span class=\"text-muted-foreground text-sm\">© 2025 Ikenna Okpala.</span><div class=\"flex flex-row gap-3\"><a href=\"https://github.com/Ikenna-Okpala\"><uk-icon icon=\"github\" width=\"20\" height=\"20\" cls-custom=\"text-muted-foreground uk-transition-scale-up uk-transition-opaque\" class=\"uk-transition-toggle\"></uk-icon></a> <a href=\"https://www.linkedin.com/in/ikennajesse/\"><uk-icon icon=\"linkedin\" width=\"20\" height=\"20\" cls-custom=\"text-muted-foreground uk-transition-scale-up uk-transition-opaque\" class=\"uk-transition-toggle\"></uk-icon></a> <a href=\"https://www.instagram.com/ikenna_jess_/\"><uk-icon icon=\"instagram\" width=\"20\" height=\"20\" cls-custom=\"text-muted-foreground uk-transition-scale-up uk-transition-opaque\" class=\"uk-transition-toggle\"></uk-icon></a></div></footer>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func App(title string, content templ.Component) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<html lang=\"en\"><head><script>\r\n\t\r\n  const htmlElement = document.documentElement;\r\n\r\n  if (\r\n    localStorage.getItem(\"mode\") === \"dark\" ||\r\n    (!(\"mode\" in localStorage) &&\r\n      window.matchMedia(\"(prefers-color-scheme: dark)\").matches)\r\n  ) {\r\n    htmlElement.classList.add(\"dark\");\r\n  } else {\r\n    htmlElement.classList.remove(\"dark\");\r\n  }\r\n\r\n  htmlElement.classList.add(\r\n    localStorage.getItem(\"theme\") || \"uk-theme-emerald\",\r\n  );\r\n</script><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/view/app.templ`, Line: 51, Col: 17}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</title><link rel=\"stylesheet\" href=\"/static/css/output.css\"><link rel=\"stylesheet\" href=\"/static/css/prism-theme.css\"><link rel=\"icon\" href=\"./favicon.ico\" type=\"image/x-icon\"><link rel=\"preconnect\" href=\"https://fonts.googleapis.com\"><link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin><link href=\"https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&amp;display=swap\" rel=\"stylesheet\"><link rel=\"stylesheet\" href=\"https://unpkg.com/franken-ui@2.0.0-internal.45/dist/css/core.min.css\"><script src=\"https://cdnjs.cloudflare.com/ajax/libs/prism/1.23.0/prism.min.js\"></script><script src=\"https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/components/prism-jsx.min.js\"></script><script src=\"https://unpkg.com/franken-ui@2.0.0-internal.45/dist/js/core.iife.js\" type=\"module\"></script><script src=\"https://unpkg.com/franken-ui@2.0.0-internal.45/dist/js/icon.iife.js\" type=\"module\"></script><script src=\"https://unpkg.com/htmx.org@2.0.4\" integrity=\"sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+\" crossorigin=\"anonymous\"></script></head><script>\r\n\t\tvar socket = new WebSocket(\"ws://localhost:8080/live\")\r\n\r\n\tsocket.onclose = (event) => {\r\n\t\tconsole.log(\"The conn has been closed....\")\r\n\t\t\r\n\t\tsetTimeout(() => {\r\n\t\t\tlocation.reload()\r\n\t\t}, 3000)\r\n\t}\r\n\t\t</script><body class=\"bg-background text-foreground\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Nav().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "<div id=\"tab-switcher\" class=\"flex flex-col pt-28 px-32 pb-20 w-full justify-center\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = content.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Footer().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
