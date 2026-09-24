import { defineConfig, presetWebFonts, presetUno } from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetWebFonts({
      provider: 'none',
      fonts: {
        sans: 'Inter, ui-sans-serif, system-ui',
        mono: 'Fira Code, ui-monospace, SFMono-Regular',
      },
    }),
  ],
})
