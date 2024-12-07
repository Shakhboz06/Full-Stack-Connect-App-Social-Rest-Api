// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: ['@nuxt/ui'],
  devServer:{
    port: 3000
  },
  app:{
    pageTransition: {name: "fade", mode: "out-in" }
  },
  runtimeConfig:{
    public: {
      BASE_URL: process.env.BACKEND_URL,
    }
  }
})