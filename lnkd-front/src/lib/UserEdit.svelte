<script>
    import {modalIdStore, LargeModal, ModalButton} from 'flowbite-svelte'
    import { afterUpdate } from 'svelte';
    import CryptoJS from "crypto-js";
    export let openModel = false;
    export let editUser = {};
    export let usersToken = '';

    let userPassword = '';
    let confirmPassword = '';
    let userEmail = '';

    let inputClass = 'block py-2.5 px-0 w-full text-sm text-gray-900 bg-transparent border-0 border-b-2 border-gray-300 appearance-none dark:text-white dark:border-gray-600 dark:focus:border-blue-500 focus:outline-none focus:ring-0 focus:border-blue-600 peer';
    let labelClass = 'absolute text-sm text-gray-500 dark:text-gray-400 duration-300 transform -translate-y-6 scale-75 top-3 -z-10 origin-[0] peer-focus:left-0 peer-focus:text-blue-600 peer-focus:dark:text-blue-500 peer-placeholder-shown:scale-100 peer-placeholder-shown:translate-y-0 peer-focus:scale-75 peer-focus:-translate-y-6';


    afterUpdate(() => {
        if (editUser.hasOwnProperty('email') && editUser.email) {
            userEmail = editUser.email;
        }
        if (openModel) {
            modalIdStore.update((n) => (n = "user-edit-modal"));
        } else {
            closeModal()
        }
    })

    export const closeModal = () => {
        modalIdStore.update((value) => (value = null));
    };

    export const handleSubmit = () => {
        let userData = {
            password: CryptoJS.SHA256(userPassword).toString(),
            username: editUser.username,
            email: userEmail,
        }
        fetch('/api/user', {
            method: 'PATCH',
            headers: {
                'Authorization': usersToken,
            },
            body: JSON.stringify(userData),
        })
            .then(response => response.json())
            .then(data => {
                closeModal();
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    };
</script>

<LargeModal
        id="user-edit-modal"
        btnColor="indigo"
        title="User edit"
        btn1="Submit"
        btn2="Cancel"
        on:handlebtn1={handleSubmit}
        on:handlebtn2={closeModal}
>
    <div class="user-edit-name-field">
        Editing: {editUser.username}
    </div>

    <div class="relative z-0 mb-6 w-full group">
        <input id="new-user-email" bind:value={userEmail} type="email" class={inputClass} required
               placeholder=" "/>
        <label for="new-user-email" class={labelClass}>Email</label>
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

</LargeModal>