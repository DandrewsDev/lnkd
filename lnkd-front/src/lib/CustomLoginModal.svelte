<script lang="ts">
    import {modalIdStore} from "flowbite-svelte";
    import CryptoJS from 'crypto-js';
    import axios from 'axios';
    import {userJwt} from "../stores.js"

    export let id: string = 'signin-modal';
    export let titleSignIn: string = 'Sign in to our platform';
    export let lostPasswordLink: string;
    export let rememberMe = false;
    export let emailLogin = false;
    export let signUp: string;
    export let emailPlaceholder: string = 'name@company.com';
    export let usernamePlaceholder: string = 'admin';
    export const closeModal = () => {
        modalIdStore.update((value) => (value = null));
    };

    let emailUsername
    let password

    export const handleSubmit = () => {
        let loginData = {
            password: CryptoJS.SHA256(password).toString()
        }
        if (emailLogin) {
            loginData['email'] = emailUsername
        } else {
            loginData['username'] = emailUsername
        }

        axios.post('/login', loginData)
            .then(function (response) {
                userJwt.set(response.data.user);
                emailUsername = ""
                password = ""
                closeModal();
            })
            .catch(function (error) {
                console.log(error);
            });
    };
    let showModalId: string;
    modalIdStore.subscribe((value) => {
        showModalId = value;
    });
    let submitClass = 'w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800';
</script>

{#if showModalId === id}
    <div
            on:click={closeModal}
            class="bg-gray-900 bg-opacity-50 dark:bg-opacity-80 fixed inset-0 z-40 w-full h-full overflow-auto"
    >
        <!-- Main modal -->
        <div
                on:click|stopPropagation={() => {}}
                role="dialog"
                aria-modal="true"
                class="mx-auto my-20 px-4 w-full max-w-md h-full md:h-auto"
        >
            <div class="bg-white rounded-lg shadow dark:bg-gray-700">
                <div class="flex justify-end p-2">
                    <button type="button"
                            class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center dark:hover:bg-gray-800 dark:hover:text-white"
                            on:click={closeModal}>
                        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                            <path fill-rule="evenodd"
                                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                                  clip-rule="evenodd"/>
                        </svg>
                    </button>
                </div>
                <form on:submit|preventDefault={handleSubmit}>

                <div class="px-6 pb-4 space-y-6 lg:px-8 sm:pb-6 xl:pb-8">
                    <h3 class="text-xl font-medium text-gray-900 dark:text-white">
                        {titleSignIn}
                    </h3>
                    {#if emailLogin}
                        <div>
                            <label for="email-{id}"
                                   class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300">Your
                                email</label>
                            <input type="email" bind:value={emailUsername} name="email" id="email-{id}"
                                   class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
                                   placeholder={emailPlaceholder} required/>
                        </div>
                    {/if}
                    {#if !emailLogin}
                        <div>
                            <label for="username-{id}"
                                   class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300">Username</label>
                            <input type="text" bind:value={emailUsername} name="username" id="username-{id}"
                                   class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
                                   placeholder={usernamePlaceholder} required/>
                        </div>
                    {/if}
                    <div>
                        <label for="password-{id}"
                               class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300">Your
                            password</label>
                        <input type="password" bind:value={password} name="password" id="password-{id}"
                               placeholder="••••••••"
                               class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
                               required/>
                    </div>
                    <div class="flex justify-between">
                        {#if rememberMe}
                            <div class="flex items-start">
                                <div class="flex items-center h-5">
                                    <input id="remember-{id}" aria-describedby="remember" type="checkbox"
                                           class="w-4 h-4 bg-gray-50 rounded border border-gray-300 focus:ring-3 focus:ring-blue-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-blue-600 dark:ring-offset-gray-800"/>
                                </div>
                                <div class="ml-3 text-sm">
                                    <label for="remember" class="font-medium text-gray-900 dark:text-gray-300">Remember
                                        me</label>
                                </div>
                            </div>
                        {/if}
                        {#if lostPasswordLink}
                            <a href={lostPasswordLink} rel="external"
                               class="text-sm text-blue-700 hover:underline dark:text-blue-500">
                                <button type="button" on:click={closeModal}>Lost Password?</button>
                            </a>
                        {/if}
                    </div>
                    <button type="submit" class={submitClass}>Log in</button>
                    {#if signUp}
                        <div class="text-sm font-medium text-gray-500 dark:text-gray-300">
                            Not registered? <a href={signUp} rel="external"
                                               class="text-blue-700 hover:underline dark:text-blue-500">
                            <button type="button" on:click={closeModal}>Create account</button>
                        </a>
                        </div>
                    {/if}
                </div>
                </form>
            </div>
        </div>
    </div>
{/if}