<template>
    <div class="flex justify-center items-center w-4/6 h-screen max-w-fit ml-auto mr-auto">
        <Modal v-model="isOpen" :message="isSuccess" />
        <div class="flex column flex-col justify-center items-center w-full bg-white rounded-md">
            <div class="w-full flex justify-center ">
                <img class="h-32 w-36" src="/images/connect app logo.png" alt="">
            </div>
            <div class="flex justify-center items-center w-2/4 h-28 bg-darkblue w-full">
                <img class="h-18 w-30" src="/images/mail.png" alt="">
            </div>
            <div class="flex flex-col items-center mt-5 gap-6 w-full p-5">
                <h1 class="font-bold text-xl text-center">Confirm your email</h1>
                <p class="opacity-70">Hello, Thanks for your registration. Please click on the button to complete
                    confirmation
                    process.</p>
                <UButton @click="read"
                    class="hover:bg-dakrblue2 bg-darkblue w-96 h-12 text-lg flex justify-center mt-5">Confirm
                    my account</UButton>
                <div class="flex flex-col items-start w-full mt-5 border-b-2 border-gray pb-5">
                    <p class="opacity-70">Best regards,</p>
                    <span class="opacity-70">Your ConnectApp Team</span>
                </div>
            </div>
            <div class="socials">
                <img src="" alt="">
                <img src="" alt="">
                <img src="" alt="">
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>

let isOpen = ref<boolean>(false)
const isSuccess = ref<boolean>(false)

const read = async () => {
    isOpen.value = true

    const route = useRoute()
    const token = computed(() => route.params.token)
    
    const response = await fetch(`/api/${token.value}`, {
        method: 'PUT',
    })

    if (response.ok) {
        isSuccess.value = true
        setTimeout(() => {
            navigateTo("/")
        }, 1500);
    } else {
        isSuccess.value = false
        setTimeout(() => {
            navigateTo("/")
        }, 2500);
    }
}
</script>

<style></style>