<script setup>
defineProps({
    entry: Object
})

const show = defineModel();

</script>

<template>
   <div 
    :class="{ 'scale-100': show }"
    class="fixed top-0 left-0 w-screen h-screen flex items-center justify-center bg-gray-500 bg-opacity-50 transform scale-0 transition-transform duration-300">
    <!-- Button -->
    <div class="rounded-sm bg-white w-1/2 h-1/2 p-12 overflow-y-scroll"> 
        <!--Close modal button-->
        <button @click="show=!show" type="button" class="focus:outline-none">
            <!-- Hero icon - close button -->
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
        </button>

        <!-- Dialog content -->
        <p class="whitespace-break-spaces">
            <span class="text-lg font-bold">{{ entry.Subject }}</span>
            {{ entry.Content }}
        </p>
    </div>
    </div> 
</template>
