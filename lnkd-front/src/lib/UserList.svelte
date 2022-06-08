<script>
    import { SmallModal, Table } from 'flowbite-svelte'
    import {modalIdStore} from "flowbite-svelte";
    import { onMount } from "svelte";
    import CustomTableRow from "./CustomTableRow.svelte";
    import UserEdit from "./UserEdit.svelte";
    import CryptoJS from "crypto-js";
    let editUserOpen = false;
    let editUser = {};
    let smallModal1;

    export let usersToken = '';
    let userList = []
    let headerEx = ['Username','Email','Roles','Status', 'Edit', 'Remove']

    let userToDelete = {};

    const handleEditUser = (item) => {
        editUser.username = item.detail.row[0];
        editUser.email = item.detail.row[1];
        editUserOpen = true;
    }
    const handleRemoveUser = (item) => {
        editUserOpen = false;
        userToDelete.username = item.detail.row[0];
        modalIdStore.update((n) => (n = "confirm-user-removal"));
    }

    const confirmUserRemove = () => {
        fetch('/api/user/' + userToDelete.username, {
            method: 'DELETE',
            headers: {
                'Authorization': usersToken,
            },
        })
            .then(response => response.json())
            .then(data => {
                closeModal();
                fetchUserList();
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    }
    export const closeModal = () => {
        modalIdStore.update((value) => (value = null));
        editUserOpen = false;
    };

    const fetchUserList = () => {
        userList = [];
        fetch('/api/user/', {
            method: 'GET',
            headers: {
                'Authorization': usersToken,
            },
        }).then(response => response.json())
            .then(data => {
                data.users.forEach((user) => {
                    let userItem = [user.username, user.email, user.roles, user.status, '|Edit|']
                    if (user.username !== "admin") {
                        userItem.push('|Remove|');
                    }
                    userList.push(userItem);
                })
                userList = userList
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    }

    onMount(async () => {
        fetchUserList()
    });
</script>
<UserEdit openModel="{editUserOpen}" editUser="{editUser}" usersToken="{usersToken}" />
<SmallModal
        id="confirm-user-removal"
        btnColor="pink"
        title="Remove user"
        btn1="Confirm"
        btn2="Close"
        on:handlebtn1={confirmUserRemove}
        on:handlebtn2={closeModal}
>
    Confirm deletion of:{userToDelete.username}
</SmallModal>

<Table header={headerEx} >
    <CustomTableRow items={userList} on:handleRowActionOne={handleEditUser} on:handleRowActionTwo={handleRemoveUser} html />
</Table>