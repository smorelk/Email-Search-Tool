<script setup>
import { ref } from 'vue';
import Table from './components/Table.vue';

const results = ref(5);
const searchPhrase = ref("");
var entries = ref([]);

function fetchMails() {
    console.log("Search phrase:", searchPhrase.value);
    const api = `http://localhost:8080/search?q=${searchPhrase.value}`
    fetch(api, {
        mode: 'cors',
        method: "get"
    })
    .then(resp => resp.json())
    .then(data => {
        entries.value = data['hits']['hits'].map(x => {
            return {
                To: x["_source"]["To"],
                From: x["_source"]["From"],
                Date: x["_source"]["Date"],
                Subject: x["_source"]["Subject"],
                Cc: x["_source"]["Cc"],
                Content: x["_source"]["Content"]
            }
        });
        console.log(entries)
    })
    .catch(error => console.error(error))
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
            <!-- <div class="h-10" style="display: table-cell, vertical-align: middle;"> -->
            <!-- <a href="#" class="text-md no-underline text-black hover:text-blue-dark ml-2 px-1">Link1</a> -->
            <!-- <a href="#" class="text-md no-underline text-grey-darker hover:text-blue-dark ml-2 px-1">Link2</a> -->
            <!-- <a href="/two" class="text-lg no-underline text-grey-darkest hover:text-blue-dark ml-2">About Us</a> -->
            <!-- </div> -->
        </div>
    </nav>
</header>
<article class="bg-[url('/content-bg.jpg')] bg-cover to-cyan-900 min-h-screen">
    <!-- Search bar input -->
    <!--
    <section>
        <div class='flex items-center justify-center'>
            <div class="flex w-full mx-10 my-10 rounded">
                <input class="bg-slate-100 opacity-60 w-full px-4 py-1 text-gray-400 outline-none" type="search" name="search" placeholder="Search..." />
                    <button type="submit" class="m-2 rounded bg-blue-600 px-4 py-2 text-white">
                        <svg class="fill-current h-6 w-6" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Capa_1" x="0px" y="0px" viewBox="0 0 56.966 56.966" style="enable-background:new 0 0 56.966 56.966;" xml:space="preserve" width="512px" height="512px">
                            <path d="M55.146,51.887L41.588,37.786c3.486-4.144,5.396-9.358,5.396-14.786c0-12.682-10.318-23-23-23s-23,10.318-23,23  s10.318,23,23,23c4.761,0,9.298-1.436,13.177-4.162l13.661,14.208c0.571,0.593,1.339,0.92,2.162,0.92  c0.779,0,1.518-0.297,2.079-0.837C56.255,54.982,56.293,53.08,55.146,51.887z M23.984,6c9.374,0,17,7.626,17,17s-7.626,17-17,17  s-17-7.626-17-17S14.61,6,23.984,6z" />
                        </svg>
                    </button>
            </div>
        </div>
    </section>
-->

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
            <Table :entries="entries" v-model="results"></Table>
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
