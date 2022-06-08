<script>
    import {Badge, modalIdStore, SmallModal, Tooltip} from 'flowbite-svelte';
    import {createEventDispatcher, onMount} from 'svelte';
    import {Trash} from "svelte-heros";

    let lists = [];
    let selectedItem = {};
    let title = "Your LNKDs";
    let smallModal1;

    let containerClass = 'link-list-container'
    let listItemClass = 'py-3 sm:py-4'

    let listQuery = ''

    export let usersToken = "";
    export let adminControl = false;

    if (adminControl) {
        containerClass = 'link-list-container-admin';
        listItemClass = 'py-1';
        listQuery = '?showAll=true';
        title = "All LNKDs";
    }

    export const updateLnkdList = () => {
        lists = [];
        fetchUserLnkds();
    };

    const fetchUserLnkds = () => {
        const url = '/api/lnkd/' + listQuery;
        fetch(url, {
            method: 'GET',
            headers: {
                'Authorization': usersToken,
            },
        }).then(response => response.json())
            .then(data => {
                data.links.forEach((link) => {
                    lists.push(link);
                })
                lists = lists
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    }

    const deleteLink = (deleteItem) => {
        selectedItem = deleteItem
        modalIdStore.update((n) => (n = "confirm-delete-modal"));
    }
    onMount(async () => {
        fetchUserLnkds()
    });
    const closeModal = () => {
        modalIdStore.update((value) => (value = null));
    };
    const handleDelete = () => {
        fetch('/api/lnkd/' + selectedItem._id, {
            method: 'DELETE',
            headers: {
                'Authorization': usersToken,
            },
        }).then(response => response.json())
            .then(data => {
                closeModal()
            })
            .catch((error) => {
                console.error('Error:', error);
            });

        updateLnkdList()
    }
</script>

<SmallModal
        bind:this={smallModal1}
        id="confirm-delete-modal"
        btnColor="pink"
        title="Delete Link"
        btn1="Confirm"
        btn2="Cancel"
        on:handlebtn1={handleDelete}
        on:handlebtn2={closeModal}
>
    <p class="text-sm font-medium text-gray-900 truncate dark:text-white">
        {selectedItem.lnkd_url}
    </p>
    <p class="text-sm text-gray-500 truncate dark:text-gray-400">
        {selectedItem.redirect_url}
    </p>
</SmallModal>

<div class="link-list-container">
    <div class="p-4 bg-white rounded-lg border shadow-md sm:p-8 dark:bg-gray-800 dark:border-gray-700 link-list">
        <div class="flex justify-between items-center mb-4">
            <h3 class="text-xl font-bold leading-none text-gray-900 dark:text-white">
                {title}
            </h3>
        </div>
        <div class="flow-root">
            <ul class="divide-y divide-gray-200 dark:divide-gray-700">
                {#each lists as {_id, lnkd_url, redirect_url, hits, user}, i}
                    <li class="{listItemClass}">
                        <div class="flex items-center space-x-4">
                            <Tooltip tip="Number of times your link was used.">
                                <Badge name="{hits}" />
                            </Tooltip>
                            <div class="flex-1 min-w-0">
                                <p class="text-sm font-medium text-gray-900 truncate dark:text-white">
                                    {lnkd_url}
                                </p>
                                <p class="text-sm text-gray-500 truncate dark:text-gray-400">
                                    {redirect_url}
                                </p>
                            </div>
                            <div on:click={() => deleteLink(lists[i])}>
                                <svelte:component this={Trash} class="dark:text-pink-700" size=18 />
                            </div>
                        </div>
                    </li>
                {/each}
            </ul>
        </div>
    </div>
</div>