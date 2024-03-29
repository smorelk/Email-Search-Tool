<script setup>
import { ref, reactive } from 'vue';
import Table from './components/Table.vue';

// Reactive fields
const results = ref(5);
const searchPhrase = ref("");
const currentPage = ref(1);
const boundariesState = reactive({
    start: results.value * (currentPage.value-1),
    end: results.value * currentPage.value
})

var entries = ref([]);
var pages = 0;

function fetchMails() {
    console.log("Search phrase:", searchPhrase.value);
    const api = `http://localhost:8080/search?q=${searchPhrase.value}`
    fetch(api, {
        mode: 'cors',
        method: "get"
    })
    .then(resp => resp.json())
    .then(data => {
        var d = data['hits']['hits'].map(x => x['_source'])
        entries.value = d;
    })
    .catch(error => console.error(error))

    // Get the initial pages count. This value may change later. See 'setPages'.
    pages = Math.ceil(entries.length / results.value);
}

function updateBoundaries() {
    boundariesState.start = (results.value * (currentPage.value-1));
    boundariesState.end = results.value * currentPage.value;
}

// Binded with select's onChange function.
function setPages() {
    // Since the results per page change, set the pages accordantly.
    pages = Math.ceil(entries.length / results.value);
    updateBoundaries();
}

function setNextPage() {
    currentPage.value++;
    if (currentPage.value > pages) {
        currentPage.value = pages;
    }
    updateBoundaries();
}

function setPrevPage() {
    currentPage.value--;
    if (currentPage.value == 0)
        currentPage.value = 1;
    updateBoundaries();
}
</script>

<template>
<header>
    <nav class="font-sans flex flex-col text-center content-center sm:flex-row sm:text-left sm:justify-between py-2 px-6 bg-white shadow sm:items-baseline w-full">
        <div class="mb-2 sm:mb-0 flex flex-row">
        <div class="h-10 w-10 self-center mr-2">
            <img class="h-10 w-10 self-center" src="/email-icon.jpg" />
        </div>
        <div>
            <a href="/home" class="text-2xl no-underline text-grey-darkest hover:text-blue-dark font-sans font-bold">Email Search</a><br>
            <span class="text-xs text-grey-dark">Search for emails as fast as you can type!</span>
        </div>
        </div>

        <div class="sm:mb-0 self-center">
        </div>
    </nav>
</header>
<article class="bg-[url('/content-bg.jpg')] bg-cover to-cyan-900 min-h-screen">
    <!-- Table and Content -->
    <section class="flex 2xl:flex-row xl:flex-row lg:flex-row md:flex-row sm:flex-row items-baseline p-6 font-mono w-full h-full">
        <div class="2xl:w-full xl:w-full lg:w-full sm:w-full mx-auto px-4 sm:px-8 h-full">
        <div class="py-8">
            <div>
                <h2 class="text-2xl font-semibold leading-tight">Email entries</h2>
            </div>
            <div class="my-2 flex sm:flex-row flex-col">
                <div class="flex flex-row mb-1 sm:mb-0">
                    <div class="relative">
                        <select 
                            v-model="results"
                            @change="setPages"
                            class="appearance-none h-full rounded-l border block appearance-none w-full bg-white border-gray-400 text-gray-700 py-2 px-4 pr-8 leading-tight focus:outline-none focus:bg-white focus:border-gray-500">
                            <option value="5" >5</option>
                            <option value="10">10</option>
                            <option value="15">15</option>
                        </select>
                        <div
                            class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                            <svg class="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                                <path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" />
                            </svg>
                        </div>
                    </div>
                </div>
                <div class="block relative">
                    <span class="h-full absolute inset-y-0 left-0 flex items-center pl-2">
                        <svg viewBox="0 0 24 24" class="h-4 w-4 fill-current text-gray-500">
                            <path
                                d="M10 4a6 6 0 100 12 6 6 0 000-12zm-8 6a8 8 0 1114.32 4.906l5.387 5.387a1 1 0 01-1.414 1.414l-5.387-5.387A8 8 0 012 10z">
                            </path>
                        </svg>
                    </span>
                    <input placeholder="Search"
                        v-model="searchPhrase"
                        @keyup.enter="fetchMails"
                        class="appearance-none rounded-r rounded-l sm:rounded-l-none border border-gray-400 border-b block pl-8 pr-6 py-2 w-full bg-white text-sm placeholder-gray-400 text-gray-700 focus:bg-white focus:placeholder-gray-600 focus:text-gray-700 focus:outline-none" required />
                </div>
            </div>
            <Table :entries="entries" v-model="boundariesState"></Table>
            <div
               class="px-5 py-5 bg-white border-t flex flex-col xs:flex-row items-center xs:justify-between">
                <span class="text-xs xs:text-sm text-gray-900">
                    Showing {{ boundariesState.start }} to {{ boundariesState.end-1 }} of {{ entries.length }} Entries
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
    </section>
</article>
<footer class="footer w-full p-4 bg-slate-100 text-gray-800">
    <div class="text-center">
        <p>
          Copyright © 2024 -
          <a class="font-semibold" href="mailto:ken.morel.santana@gmail.com"
            >Ken Morel </a>
        </p>
      </div>
</footer>
</template>
<style scoped>
</style>
