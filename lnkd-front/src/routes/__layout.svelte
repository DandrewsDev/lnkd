<script>
	import '../app.css';
    import { userJwt } from "../stores.js"
    import {DarkMode} from "flowbite-svelte";
    import { onDestroy } from 'svelte'
    import Nav from "../lib/Nav.svelte";

    let btnClass="text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 rounded-lg text-sm p-2.5 fixed z-50 dark-mode-footer"
    let usersToken;
    let isLoggedIn = false;
    let username = '';
    let isAdmin

    const unSubJwt = userJwt.subscribe((value) => {
        usersToken = value;
        if (value) {
            isLoggedIn = true;
            usersToken = value;
            if(typeof window !== "undefined") {
                username = window.localStorage.getItem("username")
                let userRoles = window.localStorage.getItem("roles");
                // Note this just handles the display portion. Actual access controls are implemented in the API.
                if (userRoles.includes('admin')) {
                    isAdmin = true;
                }
            }
        }
        else {
            isLoggedIn = false;
        }
    })
    onDestroy(unSubJwt)
</script>
<Nav loggedIn="{isLoggedIn}" userName="{username}" isAdmin="{isAdmin}" />

<slot />

<DarkMode {btnClass} />