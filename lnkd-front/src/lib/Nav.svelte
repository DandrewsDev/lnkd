<script>
    import {modalIdStore, Navbar, NavBrand, NavHamburger, NavLi, NavUl} from "flowbite-svelte";
    import CustomLoginModal from "./CustomLoginModal.svelte";
    import {userJwt} from "../stores.js";

    export let loggedIn = false;
    export let isAdmin = false;
    export let userName = '';

    function handleLoginEvent() {
        modalIdStore.update((n) => (n = "signin1"));
    }

    function handleLogoutEvent() {
        loggedIn = false;
        userJwt.set('');
        window.location.href = '/';
    }
</script>

<Navbar let:hidden let:toggle>
    <NavBrand href="/">
        <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
			LNKD
		</span>
    </NavBrand>
    <NavHamburger on:click={toggle}/>
    <NavUl {hidden}>
        {#if !loggedIn}
            <div
                    class="block py-2 pr-4 pl-3 text-gray-700 border-b border-gray-100 hover:bg-gray-50 md:hover:bg-transparent md:border-0 md:hover:text-blue-700 md:p-0 dark:text-gray-400 md:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent dark:border-gray-700"
                    on:click={handleLoginEvent}>
                Login
            </div>
        {/if}
        {#if loggedIn}
            <span class="self-center whitespace-nowrap dark:text-white">
			    Welcome user: {userName}
		    </span>
            {#if isAdmin}
                <NavLi href="/UserManagement">User Management</NavLi>
                <NavLi href="/RouteManagement">Route Management</NavLi>
            {/if}
            <div
                    class="block py-2 pr-4 pl-3 text-gray-700 border-b border-gray-100 hover:bg-gray-50 md:hover:bg-transparent md:border-0 md:hover:text-blue-700 md:p-0 dark:text-gray-400 md:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent dark:border-gray-700"
                    on:click={handleLogoutEvent}>
                Logout
            </div>
        {/if}
    </NavUl>
</Navbar>

<CustomLoginModal
        id="signin1"
        titleSignIn="Log in"
/>