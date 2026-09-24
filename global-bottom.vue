<template>
  <footer
    v-if="$nav.currentLayout !== 'cover'"
    class="absolute bottom-0 right-0 p-2 text-sm text-gray-400 z-50"
  >
    {{ $nav.currentPage }} / {{ $nav.total }}
  </footer>
</template>

<style>
/* Force inline code to be readable on white backgrounds */
.slidev-layout :not(pre) > code {
  color: #1a5c5c !important;
  background-color: #e6f4f1 !important;
}
.slidev-layout :not(pre) > code span {
  color: inherit !important;
}

/* Fix Slidev v52+ Presenter Mode Layout Overflow */
.grid-container {
  width: 100vw !important;
}
</style>
