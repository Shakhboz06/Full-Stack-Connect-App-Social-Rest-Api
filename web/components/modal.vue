<template>
    <div>
        <UModal v-model="isOpen">
            <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
                <template #header>
                    <div class="flex items-center justify-between">
                        <h3 class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
                            User Confirmation
                        </h3>
                    </div>
                </template>
                <h1 v-if="message">Congratulation you successfully activated your account.</h1>
                <h1 v-else>We are sorry, something went wrong. Try again!</h1>
                <div class="flex justify-center items-center">
                    <UButton class="hover:bg-dakrblue2 bg-darkblue w-48 items-center h-10 flex justify-center mt-5">
                        Return to homepage
                    </UButton>
                </div>
            </UCard>
        </UModal>
    </div>
</template>

<script lang="ts" setup>
const props = defineProps({
    message: Boolean,
    modelValue : Boolean
})


const emit = defineEmits(['update:modelValue', 'update:message'])
const isOpen = computed({
    get: () => props.modelValue,
    set: (value) =>{
        emit('update:modelValue',value)
    }
})


</script>

<style></style>