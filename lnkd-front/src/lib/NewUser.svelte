<script>
    import {Alert, Button, Spinner} from "flowbite-svelte";
    import CryptoJS from 'crypto-js';

    let isBusy = false;
    let showModal = false;
    let invalidUsername = false;
    let invalidPassword = false;
    let title = 'Add new user';

    export let usersToken = '';
    const alphaOnly = /[^a-zA-Z0-9]/gm;

    let userName = '';
    let userEmail = '';
    let userPassword = '';
    let confirmPassword = '';

    let inputClass = 'block py-2.5 px-0 w-full text-sm text-gray-900 bg-transparent border-0 border-b-2 border-gray-300 appearance-none dark:text-white dark:border-gray-600 dark:focus:border-blue-500 focus:outline-none focus:ring-0 focus:border-blue-600 peer';
    let labelClass = 'absolute text-sm text-gray-500 dark:text-gray-400 duration-300 transform -translate-y-6 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-blue-600 peer-focus:dark:text-blue-500 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-6';

    const addNewUserHandler = () => {
        showModal = true;
    }
    const closeModal = () => {
        showModal = false;
    }
    const formValidation = () => {
        invalidUsername = alphaOnly.test(userName);
        invalidPassword = userPassword !== confirmPassword;

        return invalidUsername === false && invalidPassword === false;
    }
    const submitNewUser = () => {
        let validForm = formValidation();
        if (!validForm) {
            return
        }
        isBusy = true;
        let userData = {
            password: CryptoJS.SHA256(userPassword).toString(),
            username: userName,
            email: userEmail,
        }
        fetch('/api/user', {
            method: 'POST',
            headers: {
                'Authorization': usersToken,
            },
            body: JSON.stringify(userData),
        })
            .then(response => response.json())
            .then(data => {
                closeModal();
                isBusy = false;
            })
            .catch((error) => {
                console.error('Error:', error);
                isBusy = false;
            });
    }

    let divClass = 'relative bg-white rounded-lg shadow dark:bg-gray-700';
    let headDivClass = 'flex justify-between items-center p-5 rounded-t border-b dark:border-gray-600';
    let h3Class = 'text-xl font-medium text-gray-900 dark:text-white';
    let buttonClass = 'text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center dark:hover:bg-gray-600 dark:hover:text-white';
    let pClass = 'text-base leading-relaxed text-gray-500 dark:text-gray-400';
    let footerClass = 'flex items-center p-6 space-x-2 rounded-b border-t border-gray-200 dark:border-gray-600';
    let btn2Class = 'dark:bg-red-500 text-white text-center font-medium rounded-lg text-sm py-2.5 px-5 bg-red-400  cursor-not-allowed items-center inline-flex';

</script>
<Button class="new-user-submit" on:click={addNewUserHandler}>Add new user</Button>
{#if showModal}
    <!-- Large Modal -->
    <div
            on:click={closeModal}
            class="bg-gray-900 bg-opacity-50 dark:bg-opacity-80 fixed inset-0 z-40 w-full h-full overflow-auto"
    >
        <div
                on:click|stopPropagation={() => {}}
                role="dialog"
                aria-modal="true"
                class="mx-auto my-20 px-4 w-full max-w-4xl h-full md:h-auto"
        >
            <!-- Modal content -->
            <div class={divClass}>
                <!-- Modal header -->
                <div class={headDivClass}>
                    <h3 class={h3Class}>
                        {title}
                    </h3>
                    <button type="button" class={buttonClass} on:click={closeModal}>
                        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                            <path fill-rule="evenodd"
                                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                                  clip-rule="evenodd"/>
                        </svg>
                    </button>
                </div>
                <form on:submit|preventDefault={submitNewUser}>
                    {#if invalidUsername}
                        <Alert id="alert-red" color="red">
                                <span slot="content">
                                    <span class="font-medium">
                                        Invalid Username
                                    </span>
                                        Username can not contain special characters.
                                </span>
                        </Alert>
                    {/if}
                    {#if invalidPassword}
                        <Alert id="alert-red" color="red">
                                <span slot="content">
                                    <span class="font-medium">
                                        Invalid Password
                                    </span>
                                        Passwords do not match.
                                </span>
                        </Alert>
                    {/if}
                    <!-- Modal body -->
                    <div class="p-6 space-y-6">

                        <div class="relative z-0 mb-6 w-full group">
                            <input id="new-user-email" bind:value={userEmail} type="email" class={inputClass} required
                                   placeholder=" "/>
                            <label for="new-user-email" class={labelClass}>Email</label>
                        </div>
                        <div class="relative z-0 mb-6 w-full group">
                            <input id="new-user-username" bind:value={userName} type="text"
                                   class={inputClass} required placeholder=" "/>
                            <label for="new-user-username" class={labelClass}>Username</label>
                        </div>
                        <div class="relative z-0 mb-6 w-full group">
                            <input id="new-user-password" bind:value={userPassword} type="password" class={inputClass}
                                   required placeholder=" "/>
                            <label for="new-user-password" class={labelClass}>Password</label>
                        </div>
                        <div class="relative z-0 mb-6 w-full group">
                            <input id="new-user-confirm-password" bind:value={confirmPassword} type="password"
                                   class={inputClass} required placeholder=" "/>
                            <label for="new-user-confirm-password" class={labelClass}>Confirm Password</label>
                        </div>
                        {#if isBusy}
                            <Spinner/>
                        {/if}
                    </div>
                    <!-- Modal footer -->
                    <div class={footerClass}>
                        <Button type="submit">Add new User</Button>
                        <Button on:click={closeModal} class="{btn2Class}">Cancel</Button>
                    </div>
                </form>
            </div>
        </div>
    </div>
{/if}