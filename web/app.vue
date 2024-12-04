<template>
  <div>
    <transition name="fade">
      <NuxtLayout>
        <NuxtPage />
        <NuxtRouteAnnouncer />
      </NuxtLayout>
    </transition>
    </div>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>