<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow-sm border-b">
      <div class="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
        <h1 class="text-xl font-bold text-gray-800">📋 test-runner 게시판</h1>
        <button
          @click="showForm = !showForm; editPost = null"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors text-sm"
        >
          {{ showForm ? '닫기' : '글쓰기' }}
        </button>
      </div>
    </header>

    <main class="max-w-4xl mx-auto px-4 py-6">
      <!-- Create/Edit form -->
      <div v-if="showForm" class="bg-white rounded-xl shadow-sm border p-6 mb-6">
        <h2 class="text-lg font-semibold mb-4">{{ editPost ? '글 수정' : '새 글 작성' }}</h2>
        <div class="space-y-3">
          <input
            v-model="formTitle"
            placeholder="제목"
            class="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
          />
          <input
            v-model="formAuthor"
            placeholder="작성자 (비워두면 익명)"
            class="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none"
          />
          <textarea
            v-model="formContent"
            placeholder="내용을 입력하세요"
            rows="5"
            class="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none resize-y"
          ></textarea>
          <div class="flex gap-2">
            <button
              @click="submitPost"
              :disabled="!formTitle.trim()"
              class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
            >
              {{ editPost ? '수정 완료' : '등록' }}
            </button>
            <button
              @click="showForm = false; editPost = null"
              class="px-4 py-2 bg-gray-100 text-gray-600 rounded-lg hover:bg-gray-200 transition-colors text-sm"
            >
              취소
            </button>
          </div>
        </div>
      </div>

      <!-- Detail view -->
      <div v-if="detailPost" class="bg-white rounded-xl shadow-sm border p-6 mb-6">
        <div class="flex items-start justify-between mb-4">
          <div>
            <h2 class="text-xl font-bold text-gray-800">{{ detailPost.title }}</h2>
            <p class="text-sm text-gray-400 mt-1">
              {{ detailPost.author }} · {{ formatDate(detailPost.createdAt) }}
            </p>
          </div>
          <div class="flex gap-2">
            <button @click="editFromDetail" class="text-sm text-indigo-600 hover:text-indigo-800">수정</button>
            <button @click="deletePost(detailPost.id)" class="text-sm text-red-500 hover:text-red-700">삭제</button>
            <button @click="detailPost = null" class="text-sm text-gray-400 hover:text-gray-600">닫기</button>
          </div>
        </div>
        <div class="prose prose-sm max-w-none text-gray-700 whitespace-pre-wrap">{{ detailPost.content }}</div>
      </div>

      <!-- Post list -->
      <div class="space-y-3">
        <div
          v-for="post in posts"
          :key="post.id"
          @click="viewPost(post)"
          class="bg-white rounded-xl shadow-sm border p-5 cursor-pointer hover:border-indigo-300 hover:shadow transition-all"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-gray-800 truncate">{{ post.title }}</h3>
              <p class="text-sm text-gray-400 mt-1">
                {{ post.author }} · {{ formatDate(post.createdAt) }}
              </p>
            </div>
            <div class="flex gap-2 ml-4 shrink-0">
              <button
                @click.stop="editPostClick(post)"
                class="text-xs text-indigo-500 hover:text-indigo-700"
              >수정</button>
              <button
                @click.stop="deletePost(post.id)"
                class="text-xs text-red-400 hover:text-red-600"
              >삭제</button>
            </div>
          </div>
          <p class="text-sm text-gray-500 mt-2 line-clamp-2">{{ post.content }}</p>
        </div>

        <div v-if="posts.length === 0 && !loading" class="text-center py-12 text-gray-400">
          게시글이 없습니다. 첫 글을 작성해보세요!
        </div>
        <div v-if="loading" class="text-center py-12 text-gray-400">로딩 중...</div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const posts = ref([])
const loading = ref(true)
const showForm = ref(false)
const editPost = ref(null)
const detailPost = ref(null)
const formTitle = ref('')
const formContent = ref('')
const formAuthor = ref('')

async function fetchPosts() {
  try {
    const res = await fetch('/api/posts')
    posts.value = await res.json()
  } catch (e) {
    console.error('Failed to fetch posts:', e)
  } finally {
    loading.value = false
  }
}

function formatDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

async function submitPost() {
  if (!formTitle.value.trim()) return

  if (editPost.value) {
    const res = await fetch(`/api/posts/${editPost.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: formTitle.value, content: formContent.value })
    })
    if (res.ok) {
      await fetchPosts()
      if (detailPost.value) {
        detailPost.value = await res.json()
      }
    }
  } else {
    const res = await fetch('/api/posts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: formTitle.value,
        content: formContent.value,
        author: formAuthor.value || '익명'
      })
    })
    if (res.ok) {
      await fetchPosts()
    }
  }

  showForm.value = false
  editPost.value = null
  formTitle.value = ''
  formContent.value = ''
  formAuthor.value = ''
}

function editPostClick(post) {
  editPost.value = post
  formTitle.value = post.title
  formContent.value = post.content
  formAuthor.value = post.author === '익명' ? '' : post.author
  showForm.value = true
  detailPost.value = null
}

function editFromDetail() {
  if (detailPost.value) {
    editPostClick(detailPost.value)
  }
}

async function deletePost(id) {
  if (!confirm('정말 삭제하시겠습니까?')) return
  const res = await fetch(`/api/posts/${id}`, { method: 'DELETE' })
  if (res.ok) {
    if (detailPost.value?.id === id) detailPost.value = null
    await fetchPosts()
  }
}

function viewPost(post) {
  detailPost.value = post
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(fetchPosts)
</script>
