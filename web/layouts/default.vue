<template>
    <div>
        <slot />
    </div>
</template>

<script lang="ts" setup>

useHead({
    title: 'Connect App',
    link: [
        {
            rel: 'preconnect',
            href: 'https://fonts.googleapis.com'
        },
        {
            rel: 'stylesheet',
            href: 'https://fonts.googleapis.com/css2?family=Inter&display=swap',
            crossorigin: ''
        }
    ]
})
</script>

<style>
*{
    margin:0;
    padding: 0;
    box-sizing: border-box;
}
body{
    font-family: Inter;
}
</style>