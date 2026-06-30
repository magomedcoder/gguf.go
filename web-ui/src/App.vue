<script setup lang="ts">
import ChatView from '@/components/ChatView.vue'
</script>

<template>
  <div class="min-h-screen bg-[#0f0f0f] font-sans text-[#e8e8e8] antialiased">
    <chat-view />
  </div>
</template>

<style>
@import 'tailwindcss';

@theme {
  --animate-blink: blink 1s step-end infinite;

  @keyframes blink {
    50% {
      opacity: 0;
    }
  }
}
</style>
