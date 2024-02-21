<script setup>
import { reactive, ref, toRaw } from 'vue';
import Dialog from './Dialog.vue'

const props = defineProps(['entries']);
const headers = ["To", "From", "Date", "Subject", "View"];
const selectedEntry = ref({});
const show = ref(false);
const boundaries = defineModel();

function showEntryContent(entry) {
    selectedEntry.value = toRaw(entry);
    show.value = true;
}

</script>

<template>
    <div class="-mx-4 sm:-mx-8 px-4 sm:px-8 py-4">
        <div class="inline-block min-w-full shadow rounded-lg">
            <table class="min-w-full leading-normal sm:w-full md:w-full lg:w-full xl:w-full 2xl:w-full">
                <!-- Table Header -->
                <thead>
                <tr>
                    <th
                        v-for="header in headers"
                        class="border px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        {{ header }}
                     </th>
                </tr>
                </thead>
                    <tbody>
                       <!-- Table Rows -->
                       <tr v-for="entry in props.entries.slice(boundaries.start, boundaries.end)">
                                <td class="p-5 border-b border-gray-200 bg-white text-sm border text-balance">
                                   <p v-for="mail in entry.To" class="text-gray-900">
                                        {{ mail }}
                                    </p>
                                </td>
                                <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm border">
                                    <p class="text-gray-900 whitespace-no-wrap">{{ entry.From }}</p>
                                </td>
                                <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm border text-nowrap">
                                    <p class="text-gray-900">
                                        {{ entry.Date }}
                                    </p>
                                </td>
                                <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm border text-nowrap">
                                    <p class="text-gray-900">
                                        {{ entry.Subject }}
                                    </p>
                                </td>
                                <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm border">
                                    <button @click="showEntryContent(entry)" class="bg-grey-light hover:bg-grey text-grey-darkest font-bold py-2 px-4 rounded inline-flex items-center">
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" />
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                                        </svg>
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

    <!-- Dialog for the content -->
    <Dialog v-model="show" :entry="selectedEntry"></Dialog>
</template>
