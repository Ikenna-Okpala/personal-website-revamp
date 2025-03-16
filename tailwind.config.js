const plugin = require('tailwindcss/plugin')

import franken from "franken-ui/shadcn-ui/preset-quick";
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internal/view/**/*.{html,templ,go}",],
  theme: {
    extend: {},
  },
  plugins: [
    plugin(function ({ addVariant }) {
      addVariant('htmx-settling', ['&.htmx-settling', '.htmx-settling &'])
      addVariant('htmx-request', ['&.htmx-request', '.htmx-request &'])
      addVariant('htmx-swapping', ['&.htmx-swapping', '.htmx-swapping &'])
      addVariant('htmx-added', ['&.htmx-added', '.htmx-added &'])
    }),

  ],
  presets: [
    franken({
      customPalette: {

        ".uk-theme-emerald": {
          "--background": "262.1 41% 95%",
          "--foreground": "262.1 5% 10%",
          "--card": "262.1 41% 90%",
          "--card-foreground": "262.1 5% 15%",
          "--popover": "262.1 41% 95%",
          "--popover-foreground": "262.1 95% 10%",
          "--primary": "262.1 88.3% 57.8%",
          "--primary-foreground": "0 0% 100%",
          "--secondary": "262.1 30% 70%",
          "--secondary-foreground": "0 0% 0%",
          "--muted": "224.10000000000002 30% 85%",
          "--muted-foreground": "262.1 5% 35%",
          "--accent": "224.10000000000002 30% 80%",
          "--accent-foreground": "262.1 5% 15%",
          "--destructive": "0 50% 30%",
          "--destructive-foreground": "262.1 5% 90%",
          "--border": "262.1 30% 50%",
          "--input": "262.1 30% 25%",
          "--ring": "262.1 88.3% 57.8%",
          "--radius": "0.5rem"
        },
        ".dark.uk-theme-emerald": {
          "--background": "262.1 41% 10%",
          "--foreground": "262.1 5% 90%",
          "--card": "262.1 41% 10%",
          "--card-foreground": "262.1 5% 90%",
          "--popover": "262.1 41% 5%",
          "--popover-foreground": "262.1 5% 90%",
          "--primary": "262.1 88.3% 57.8%",
          "--primary-foreground": "0 0% 100%",
          "--secondary": "262.1 30% 20%",
          "--secondary-foreground": "0 0% 100%",
          "--muted": "224.10000000000002 30% 25%",
          "--muted-foreground": "262.1 5% 60%",
          "--accent": "224.10000000000002 30% 25%",
          "--accent-foreground": "262.1 5% 90%",
          "--destructive": "0 50% 30%",
          "--destructive-foreground": "262.1 5% 90%",
          "--border": "262.1 30% 25%",
          "--input": "262.1 30% 25%",
          "--ring": "262.1 88.3% 57.8%",
          "--radius": "0.5rem"
        }


      }
    }
    )
  ]
}