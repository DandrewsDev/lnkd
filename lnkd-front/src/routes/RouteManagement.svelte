<script>
    import {userJwt} from "../stores.js"
    import {onDestroy} from "svelte";
    import LinkList from "../lib/LinkList.svelte";

    let userToken = '';
    let linkList;

    const unSubJwt = userJwt.subscribe((value) => {
        if (value) {
            userToken = value;
        }
    })
    onDestroy(unSubJwt)
</script>
<div class="user-list-container">
    <LinkList usersToken="{userToken}" bind:this={linkList} adminControl="true" />
</div>