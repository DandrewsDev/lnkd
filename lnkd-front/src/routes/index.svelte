<script>
    import CreateLnkd from "../lib/CreateLnkd.svelte";
    import LinkList from "../lib/LinkList.svelte";
    import {onDestroy} from "svelte";
    import {userJwt} from "../stores.js"

    let isLoggedIn = false;
    let isAdmin = false;
    let userToken = '';
    let linkList;
    let username = '';

    const unSubJwt = userJwt.subscribe((value) => {
        if (value) {
            isLoggedIn = true;
            userToken = value;
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

    function handleMessage(event) {
        if (isLoggedIn) {
            linkList.updateLnkdList();
        }
    }
</script>

<CreateLnkd on:lnkd-created={handleMessage} />
{#if isLoggedIn}
    <LinkList usersToken="{userToken}" bind:this={linkList} />
{/if}