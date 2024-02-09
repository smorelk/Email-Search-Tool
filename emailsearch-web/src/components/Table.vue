<script setup>
import { reactive, ref, toRaw } from 'vue';
import Dialog from './Dialog.vue'

const props = defineProps(['entries']);
const headers = ["To", "From", "Date", "Subject", "Cc", "View"];
const showEntry = ref({});
const hidden = ref(false);

const start = ref(0);
const end = defineModel();
let currentPage = 1;

function showEntryContent(entry) {
    showEntry.value = toRaw(entry);
    hidden.value = true;
}

function setPrevPage() {

}

function setNextPage() {
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
                       <tr v-for="entry in props.entries.slice(start, end)">
                                <td class="p-5 border-b border-gray-200 bg-white text-sm border text-balance">
                                   <p v-for="mail in entry.To" class="text-gray-900">
                                        {{  mail }},
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
                                <td class="p-5 border-b border-gray-200 bg-white text-sm border text-balance">
                                    <p v-for="mail in entry.Cc" class="text-gray-900">
                                        {{  mail }},
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
                    <div
                        class="px-5 py-5 bg-white border-t flex flex-col xs:flex-row items-center xs:justify-between">
                        <span class="text-xs xs:text-sm text-gray-900">
                            Showing {{ start }} to {{ end * currentPage}} of {{ props.entries.length }} Entries
                        </span>
                        <div class="inline-flex mt-2 xs:mt-0">
                            <button
                                @click="setPrevPage"
                                class="text-sm bg-gray-300 hover:bg-gray-400 text-gray-800 font-semibold py-2 px-4 rounded-l">
                                Prev
                            </button>
                            <button
                                @click="setNextPage"
                                class="text-sm bg-gray-300 hover:bg-gray-400 text-gray-800 font-semibold py-2 px-4 rounded-r">
                                Next
                            </button>
                        </div>
                    </div>
                </div>
            </div>

    <!-- Dialog for the content -->
    <Dialog :hidden="hidden" :entry="showEntry"></Dialog>
</template>