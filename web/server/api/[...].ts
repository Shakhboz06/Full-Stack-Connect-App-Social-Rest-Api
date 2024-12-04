export default defineEventHandler(async event =>{

    
    const config = useRuntimeConfig()
    
    const token: string | undefined  = event.node.req.url?.split('/').pop()

    if(token == undefined){
        throw new Error('Requested token is undefined');
    }

    const response = await $fetch(config.public.BASE_URL + `/users/activation/${token}`,{
        method: "PUT",
    }) 
    
    return response
})