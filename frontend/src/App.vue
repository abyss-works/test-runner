<template>
  <div class="min-h-screen bg-stone-100">
    <!-- Header -->
    <header class="bg-white border-b border-stone-200 print:hidden">
      <div class="max-w-3xl mx-auto px-4 sm:px-6 py-5 flex items-center justify-between">
        <h1 class="text-lg font-serif font-bold text-stone-800 tracking-tight">
          test-runner
          <span class="text-stone-400 font-normal mx-1">/</span>
          <span class="text-sm font-sans font-normal text-stone-500">게시판</span>
        </h1>
        <button
          @click="openNewForm"
          class="px-4 py-1.5 text-sm bg-stone-800 text-stone-50 rounded hover:bg-stone-700 transition-colors font-serif tracking-wide"
        >
          + 새 글
        </button>
      </div>
    </header>

    <main class="max-w-3xl mx-auto px-4 sm:px-6 py-8">
      <!-- Error -->
      <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded text-sm text-red-700">
        {{ error }}
        <button @click="error = ''" class="ml-2 underline">닫기</button>
      </div>

      <!-- Detail View -->
      <article v-if="detailPost" class="mb-8 animate-fade-in">
        <div class="bg-white rounded-sm border border-stone-200 shadow-sm">
          <div class="p-6 sm:p-8">
            <div class="flex items-center justify-between mb-2">
              <button @click="detailPost = null" class="text-xs text-stone-400 hover:text-stone-600 tracking-wide uppercase font-sans">
                ← 목록으로
              </button>
              <div class="flex gap-3">
                <button @click="editFromDetail" class="text-xs text-stone-500 hover:text-stone-700 font-sans tracking-wide uppercase">수정</button>
                <button @click="deletePost(detailPost.id)" class="text-xs text-red-400 hover:text-red-600 font-sans tracking-wide uppercase">삭제</button>
              </div>
            </div>

            <h2 class="text-2xl font-serif font-bold text-stone-800 mt-4 mb-2 leading-snug">
              {{ detailPost.title }}
            </h2>
            <p class="text-sm text-stone-400 font-sans mb-6">
              {{ detailPost.author }} · {{ formatDate(detailPost.createdAt) }}
            </p>
            <div class="prose prose-stone prose-sm max-w-none font-serif text-stone-700 leading-relaxed whitespace-pre-wrap border-t border-stone-100 pt-6">
              {{ detailPost.content }}
            </div>

            <!-- Comments -->
            <div class="border-t border-stone-100 mt-8 pt-6">
              <h3 class="text-sm font-sans font-semibold text-stone-500 uppercase tracking-wide mb-4">
                댓글 {{ comments.length }}개
              </h3>

              <div class="space-y-3 mb-6" v-if="comments.length > 0">
                <div v-for="c in comments" :key="c.id"
                  class="flex items-start gap-3 py-3 border-b border-stone-50 last:border-0">
                  <div class="w-6 h-6 rounded-full bg-stone-200 flex items-center justify-center text-[10px] text-stone-500 font-sans shrink-0 mt-0.5">
                    {{ c.author[0] }}
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2">
                      <span class="text-xs font-sans font-medium text-stone-600">{{ c.author }}</span>
                      <span class="text-[10px] text-stone-300 font-sans">{{ formatDate(c.createdAt) }}</span>
                      <button
                        @click.stop="deleteComment(c.id)"
                        class="ml-auto text-[10px] text-red-200 hover:text-red-400 font-sans"
                      >삭제</button>
                    </div>
                    <p class="text-sm font-serif text-stone-700 mt-1 leading-relaxed">{{ c.content }}</p>
                  </div>
                </div>
              </div>
              <p v-else class="text-xs text-stone-300 font-sans mb-4">첫 댓글을 남겨보세요</p>

              <!-- Comment form -->
              <div class="flex gap-2">
                <input
                  v-model="commentAuthor"
                  placeholder="작성자"
                  class="w-24 px-2 py-1.5 text-xs bg-transparent border border-stone-200 rounded focus:border-stone-400 focus:ring-0 outline-none text-stone-600 font-sans transition-colors placeholder:text-stone-300"
                />
                <input
                  v-model="commentContent"
                  placeholder="댓글을 입력하세요..."
                  class="flex-1 px-3 py-1.5 text-sm bg-transparent border border-stone-200 rounded focus:border-stone-400 focus:ring-0 outline-none text-stone-700 font-serif transition-colors placeholder:text-stone-300"
                  @keydown.enter="submitComment"
                />
                <button
                  @click="submitComment"
                  :disabled="!commentContent.trim()"
                  class="px-3 py-1.5 text-xs bg-stone-800 text-stone-50 rounded hover:bg-stone-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors font-sans tracking-wide uppercase"
                >등록</button>
              </div>
            </div>
          </div>
        </div>
      </article>

      <!-- Create/Edit Form -->
      <div v-if="showForm" class="mb-8 animate-fade-in">
        <div class="bg-white rounded-sm border border-stone-200 shadow-sm p-6 sm:p-8">
          <h2 class="text-lg font-serif font-bold text-stone-800 mb-6">
            {{ editPost ? '글 수정' : '새 글' }}
          </h2>
          <div class="space-y-4">
            <div>
              <label class="block text-xs font-sans text-stone-500 uppercase tracking-wide mb-1">제목</label>
              <input
                v-model="formTitle"
                placeholder="제목을 입력하세요"
                class="w-full px-0 py-2 bg-transparent border-0 border-b border-stone-200 focus:border-stone-800 focus:ring-0 outline-none text-stone-800 font-serif text-base transition-colors placeholder:text-stone-300"
                @keydown.enter="submitPost"
              />
            </div>
            <div class="flex gap-4">
              <div class="flex-1">
                <label class="block text-xs font-sans text-stone-500 uppercase tracking-wide mb-1">작성자</label>
                <input
                  v-model="formAuthor"
                  placeholder="익명"
                  class="w-full px-0 py-2 bg-transparent border-0 border-b border-stone-200 focus:border-stone-800 focus:ring-0 outline-none text-stone-800 font-serif text-sm transition-colors placeholder:text-stone-300"
                />
              </div>
            </div>
            <div>
              <label class="block text-xs font-sans text-stone-500 uppercase tracking-wide mb-1">내용</label>
              <textarea
                v-model="formContent"
                placeholder="내용을 입력하세요..."
                rows="6"
                class="w-full px-0 py-2 bg-transparent border-0 border-b border-stone-200 focus:border-stone-800 focus:ring-0 outline-none text-stone-800 font-serif text-base leading-relaxed transition-colors resize-y placeholder:text-stone-300"
              ></textarea>
            </div>
            <div class="flex gap-3 pt-2">
              <button
                @click="submitPost"
                :disabled="!formTitle.trim()"
                class="px-5 py-2 bg-stone-800 text-stone-50 rounded hover:bg-stone-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors text-sm font-serif tracking-wide"
              >
                {{ editPost ? '수정' : '등록' }}
              </button>
              <button
                @click="closeForm"
                class="px-5 py-2 bg-stone-100 text-stone-600 rounded hover:bg-stone-200 transition-colors text-sm font-sans"
              >
                취소
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Post List -->
      <div class="space-y-px bg-stone-200 border border-stone-200 rounded-sm overflow-hidden" v-if="!detailPost">
        <div
          v-for="post in posts"
          :key="post.id"
          @click="viewPost(post)"
          class="bg-white hover:bg-stone-50 transition-colors cursor-pointer"
        >
          <div class="px-5 sm:px-6 py-4 flex items-start gap-4">
            <span class="text-xs text-stone-300 font-mono mt-1 w-6 shrink-0 text-right">#{{ post.id }}</span>
            <div class="flex-1 min-w-0">
              <h3 class="font-serif font-semibold text-stone-800 truncate leading-snug">
                {{ post.title }}
              </h3>
              <p class="text-xs text-stone-400 font-sans mt-1">
                {{ post.author }} · {{ formatDate(post.createdAt) }}
              </p>
              <p class="text-sm text-stone-500 font-serif mt-1.5 line-clamp-2 leading-relaxed">
                {{ post.content }}
              </p>
            </div>
            <div class="flex gap-2 shrink-0 mt-0.5" @click.stop>
              <button
                @click.stop="editPostClick(post)"
                class="text-[11px] text-stone-400 hover:text-stone-600 font-sans tracking-wide uppercase"
              >수정</button>
              <button
                @click.stop="deletePost(post.id)"
                class="text-[11px] text-red-300 hover:text-red-500 font-sans tracking-wide uppercase"
              >삭제</button>
            </div>
          </div>
        </div>

        <div v-if="posts.length === 0 && !loading" class="bg-white py-24 text-center">
          <p class="text-stone-400 font-serif text-base">게시글이 없습니다</p>
          <p class="text-stone-300 text-sm font-sans mt-2">첫 글을 남겨보세요</p>
        </div>
        <div v-if="loading" class="bg-white py-24 text-center">
          <p class="text-stone-400 font-serif text-sm">불러오는 중...</p>
        </div>
      </div>

      <!-- Footer -->
      <p class="text-center text-xs text-stone-300 font-sans mt-8 mb-4 print:hidden">
        {{ posts.length }}개의 글 · in-memory 게시판
      </p>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const posts = ref([])
const comments = ref([])
const loading = ref(true)
const error = ref('')
const showForm = ref(false)
const editPost = ref(null)
const detailPost = ref(null)
const formTitle = ref('')
const formContent = ref('')
const formAuthor = ref('')
const commentContent = ref('')
const commentAuthor = ref('')

async function fetchPosts() {
  try {
    const res = await fetch('/api/posts')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    posts.value = await res.json()
  } catch (e) {
    error.value = '게시글을 불러올 수 없습니다'
    console.error(e)
  } finally {
    loading.value = false
  }
}

function formatDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const min = String(d.getMinutes()).padStart(2, '0')
  return `${y}.${m}.${day} ${h}:${min}`
}

function resetForm() {
  formTitle.value = ''
  formContent.value = ''
  formAuthor.value = ''
  editPost.value = null
  showForm.value = false
}

function openNewForm() {
  resetForm()
  showForm.value = true
  detailPost.value = null
}

function closeForm() {
  resetForm()
}

async function submitPost() {
  if (!formTitle.value.trim()) return

  try {
    let res
    if (editPost.value) {
      res = await fetch(`/api/posts/${editPost.value.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: formTitle.value, content: formContent.value })
      })
      if (!res.ok) throw new Error('수정 실패')
      const updated = await res.json()
      if (detailPost.value?.id === updated.id) {
        detailPost.value = updated
      }
    } else {
      res = await fetch('/api/posts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: formTitle.value,
          content: formContent.value,
          author: formAuthor.value || '익명'
        })
      })
      if (!res.ok) throw new Error('등록 실패')
    }
    await fetchPosts()
    resetForm()
  } catch (e) {
    error.value = e.message
  }
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
  if (detailPost.value) editPostClick(detailPost.value)
}

async function deletePost(id) {
  if (!confirm('정말 삭제하시겠습니까?')) return
  try {
    const res = await fetch(`/api/posts/${id}`, { method: 'DELETE' })
    if (!res.ok) throw new Error('삭제 실패')
    if (detailPost.value?.id === id) detailPost.value = null
    await fetchPosts()
  } catch (e) {
    error.value = e.message
  }
}

async function viewPost(post) {
  detailPost.value = post
  window.scrollTo({ top: 0, behavior: 'smooth' })
  await fetchComments(post.id)
}

async function fetchComments(postId) {
  try {
    const res = await fetch(`/api/posts/${postId}/comments`)
    if (!res.ok) throw new Error('댓글 로드 실패')
    comments.value = await res.json()
  } catch (e) {
    console.error(e)
    comments.value = []
  }
}

async function submitComment() {
  if (!commentContent.value.trim()) return
  if (!detailPost.value) return
  try {
    const res = await fetch(`/api/posts/${detailPost.value.id}/comments`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        content: commentContent.value,
        author: commentAuthor.value || '익명'
      })
    })
    if (!res.ok) throw new Error('댓글 등록 실패')
    commentContent.value = ''
    commentAuthor.value = ''
    await fetchComments(detailPost.value.id)
  } catch (e) {
    error.value = e.message
  }
}

async function deleteComment(commentId) {
  if (!detailPost.value) return
  if (!confirm('댓글을 삭제하시겠습니까?')) return
  try {
    const res = await fetch(`/api/posts/${detailPost.value.id}/comments/${commentId}`, {
      method: 'DELETE'
    })
    if (!res.ok) throw new Error('댓글 삭제 실패')
    await fetchComments(detailPost.value.id)
  } catch (e) {
    error.value = e.message
  }
}

onMounted(fetchPosts)
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Noto+Serif+KR:wght@400;600;700&family=Noto+Sans+KR:wght@300;400;500&display=swap');

@tailwind base;
@tailwind components;
@tailwind utilities;

body {
  font-family: 'Noto Sans KR', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.font-serif {
  font-family: 'Noto Serif KR', Georgia, 'Times New Roman', serif;
}

.animate-fade-in {
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Global overrides */
input:focus, textarea:focus {
  outline: none;
  box-shadow: none;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
