<script>
import UserList from "../lib/UserList.svelte";
import NewUser from "../lib/NewUser.svelte";
import {userJwt} from "../stores.js"
import {onDestroy} from "svelte";

let userToken = '';

const unSubJwt = userJwt.subscribe((value) => {
    if (value) {
        userToken = value;
    }
})
onDestroy(unSubJwt)
</script>
<div class="user-list-container">
    <UserList usersToken="{userToken}" />
    <NewUser usersToken="{userToken}"  />
</div>