<script>
    import {Button, Input, Spinner, Toast} from 'flowbite-svelte'
    import {userJwt} from "../stores.js"
    import {createEventDispatcher, onDestroy} from "svelte";
    import {Link} from "svelte-heros";

    let lnkdName = ""
    let redirectLnkd = ""
    let isLoggedIn = false;
    let isBusy = false;
    let usersToken;
    let toasts = [];
    let localUrl = '';

    if(typeof window !== "undefined") {
        localUrl = window.location.origin;
    }

    const unSubJwt = userJwt.subscribe((value) => {
        isLoggedIn = !!value;
        if (value) {
            usersToken = value;
        }
    })

    const dispatch = createEventDispatcher();

    function handleLnkdAdded(lnkdObj) {
        dispatch('lnkd-created', {
            lnkd: lnkdObj
        });
    }
    onDestroy(unSubJwt)

    const submitNewLink = () => {
        isBusy = true;
        const data = { "redirect_url": redirectLnkd };

        if (lnkdName) {
            data['lnkd_url'] = lnkdName;
        }

        lnkdName = ""
        redirectLnkd = ""

        fetch('/api/lnkd', {
            method: 'POST',
            headers: {
                'Authorization': usersToken,
            },
            body: JSON.stringify(data),
        })
            .then(response => response.json())
            .then(data => {
                navigator.clipboard.writeText(data.lnkd);
                handleLnkdAdded(data.lnkd)
                toasts.push(data)
                toasts = toasts;
                isBusy = false;
            })
            .catch((error) => {
                console.error('Error:', error);
                isBusy = false;
            });
    };
</script>

<div class="new-lnkd">
    <form on:submit|preventDefault={submitNewLink}>
        {#if isLoggedIn}
            <Input
                    class="mb-6"
                    bind:value="{lnkdName}"
                    placeholder="custom-name"
                    label="LNKD Vanity URL (Optional)"
            />
        {/if}
        <Input
                bind:value="{redirectLnkd}"
                class="mb-6"
                id="large-input"
                name="size"
                size="sm:text-md"
                type="url"
                placeholder="https://example.com"
                label="URL To Shorten *"
                required
        />
        {#if !lnkdName}
            <p class="font-medium text-gray-300">
                URLs without a custom vanity url will be given an autogenerated link {lnkdName}
            </p>
        {/if}
        {#if lnkdName}
            <p class="font-medium text-gray-300">Your short url: {localUrl}/{lnkdName}</p>
        {/if}
        {#if isBusy}
            <Spinner />
        {/if}
        <Button class="new-lnkd-submit" type="submit">Submit</Button>
    </form>
</div>

<div class="toast-container">
    {#each toasts as toast, i}
        <Toast transitionType="slide" iconColor='green' divClass="custom-toast flex items-center w-full max-w-xs p-4 text-gray-500 bg-white rounded-lg shadow dark:text-gray-400 dark:bg-gray-600">
            <span slot="icon">
                <Link/>
            </span>
            <span slot="text">
                New Lnkd created: {toast.lnkd} Copied to clipboard.
            </span>
        </Toast>
    {/each}
</div>