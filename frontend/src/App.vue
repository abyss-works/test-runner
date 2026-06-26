<template>
  <div class="min-h-screen flex flex-col items-center justify-center p-8">
    <header class="text-center mb-12">
      <h1 class="text-5xl font-bold text-indigo-600 mb-4">{{PROJECT_NAME}}</h1>
      <p class="text-xl text-gray-500 max-w-xl">
        A fullstack application built with Go + Vue 3 + PostgreSQL.
      </p>
    </header>

    <main class="w-full max-w-2xl">
      <div class="bg-white rounded-xl shadow-md p-8 mb-6">
        <h2 class="text-2xl font-semibold mb-4">Health Check</h2>
        <div class="flex items-center gap-4">
          <span class="px-3 py-1 bg-green-100 text-green-800 rounded-full text-sm font-medium">
            {{ healthStatus }}
          </span>
          <button
            @click="checkHealth"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
          >
            Check API
          </button>
        </div>
      </div>

      <div class="bg-white rounded-xl shadow-md p-8">
        <h2 class="text-2xl font-semibold mb-4">Getting Started</h2>
        <ul class="space-y-3 text-gray-600">
          <li>▶ Edit <code class="bg-gray-100 px-2 py-0.5 rounded">frontend/src/App.vue</code> to customize this page</li>
          <li>▶ API endpoints are served from <code class="bg-gray-100 px-2 py-0.5 rounded">/api/*</code></li>
          <li>▶ Run <code class="bg-gray-100 px-2 py-0.5 rounded">docker build -t {{PROJECT_NAME}} .</code> to containerize</li>
        </ul>
      </div>
    </main>

    <footer class="mt-12 text-gray-400 text-sm">
      Scaffolded by idp-starter
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const healthStatus = ref('unknown')

async function checkHealth() {
  try {
    const res = await fetch('/api/health')
    const data = await res.json()
    healthStatus.value = data.status
  } catch {
    healthStatus.value = 'unreachable'
  }
}
</script>
