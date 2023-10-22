<script>
    import {
        Accordion,
        Alert,
        Button,
        Card,
        DatePicker,
        Icon,
        Input,
        Layout,
        Media,
        Modal,
        Portal,
    } from "stwui";
    import {
        createUser,
        getAllAds,
        getCookie,
        getMe,
        sendCredentials,
        submitApply,
        updateMyInfo,
        updateProfile,
    } from "$lib";
    import dayjs from "dayjs";
    import { onMount } from "svelte";
    import company from "../../static/company.svg";
    import save from "../../static/save.svg";
    import wage from "../../static/wage.svg";
    import workAddress from "../../static/workAddress.svg";
    import workTime from "../../static/workTime.svg";
    import { getSvg } from "../lib/index.js";

    let token = "";
    let loading = false;
    let userConnected = false;
    let idAvert;
    let message = "";
    let noMessage = false;
    let newEmail = "";
    let newPassword = "";
    let newName = "";
    let newSurname = "";
    let newTel = "";
    let editError = false;
    let creationFailed = false;
    let newBirthDate = new Date();
    let isEditing = false;
    let userName = "";
    let userSurname = "";
    let userBirthDate = "";
    let userEmail = "";
    let userTel = "";
    let userPassword = "";
    let errorPwdInvalid = false;

    function flushEntry() {
        token = "";
        idAvert = "";
        message = "";
        noMessage = false;
        newEmail = "";
        newPassword = "";
        newName = "";
        newSurname = "";
        newTel = "";
        editError = false;
        creationFailed = false;
        newBirthDate = new Date();
        isEditing = false;
        userPassword = "";
        errorPwdInvalid = false;
    }

    let allModalStates = {
        modalApply: false,
        modalLogIn: false,
        modalSignIn: false,
        modalEdit: false,
    };

    function showModal(modalName) {
        for (const [key, value] of Object.entries(allModalStates)) {
            if (key === modalName) {
                allModalStates[key] = true;
            } else {
                allModalStates[key] = false;
            }
        }
        console.log(allModalStates);
    }

    function hideModal(modalName) {
        for (const [key, value] of Object.entries(allModalStates)) {
            if (key === modalName) {
                allModalStates[key] = false;
            }
        }
        console.log(allModalStates);
    }

    async function tryEdit() {
        if (userPassword !== "") {
            if (userPassword.length > 6) {
                let updateSuccess = await updateProfile(
                    userEmail,
                    userName,
                    userSurname,
                    userTel,
                    dayjs(userBirthDate).format("YYYY-MM-DD"),
                    getCookie("token"),
                    userPassword
                );
                if (updateSuccess) {
                    isEditing = false;
                    hideModal("modalEdit");
                    flushEntry();
                    editError = false;
                    errorPwdInvalid = false;
                } else {
                    editError = true;
                    console.log(editError);
                }
            } else {
                errorPwdInvalid = true;
                console.log(editError);
            }
        } else {
            let updateSuccess = await updateMyInfo(
                userEmail,
                userName,
                userSurname,
                userTel,
                dayjs(userBirthDate).format("YYYY-MM-DD"),
                getCookie("token")
            );
            if (updateSuccess) {
                isEditing = false;
                hideModal("modalEdit");
                flushEntry();
            } else {
                editError = true;
                console.log(editError);
            }
        }
    }

    function hideAllModals() {
        for (let modal in allModalStates) {
            allModalStates[modal] = false;
        }
        loading = false;
        console.log(allModalStates);
    }

    async function apply(idAvert, message) {
        await submitApply(message, idAvert, getCookie("token"));
        let resp = await getAllAds(getCookie("token"), previous, true);
        allAdv = resp[0];
        flushEntry();
    }

    function editProfile() {
        showModal("modalEdit");
    }

    async function closeModalProfile() {
        hideModal("modalEdit");
        await fillUserInfo();
    }

    async function tryCreateUser() {
        let resp = await createUser(
            newEmail,
            newPassword,
            newName,
            newSurname,
            newTel,
            dayjs(newBirthDate).format("YYYY-MM-DD")
        );
        if (resp) {
            token = await sendCredentials(newEmail, newPassword);
            document.cookie = "token=" + token;
            hideModal("modalSignIn");
            creationFailed = false;
            flushEntry();
            await fillUserInfo();
        } else {
            creationFailed = true;
        }
    }

    let badTry = false;

    function logout() {
        document.cookie = "token= ; expires = Thu, 01 Jan 1970 00:00:00 GMT";
        userConnected = false;
    }

    function toggleLoading() {
        loading = !loading;
    }

    let login = "";
    let password = "";

    async function tryLogin() {
        token = await sendCredentials(login, password);
        if (token) {
            document.cookie = "token=" + token;
            userConnected = true;
            hideModal("modalLogIn");
            badTry = false;
            flushEntry();
            await fillUserInfo();
        } else {
            badTry = true;
        }
    }

    function openModalApply(id) {
        toggleLoading();
        showModal("modalApply");
        idAvert = id;
    }

    function openModalSignIn() {
        showModal("modalSignIn");
    }

    function closeModalApply() {
        toggleLoading();
        hideAllModals();
    }

    let open = "";

    function handleClick(item) {
        if (open === item) {
            open = "";
        } else {
            open = item;
        }
    }

    async function fillUserInfo() {
        userConnected = await getMe(getCookie("token"));
        userName = userConnected.Name;
        userSurname = userConnected.Surname;
        userBirthDate = userConnected.DateOfBirth;
        userEmail = userConnected.Email;
        userTel = userConnected.Phone;
        console.log(userConnected);
    }

    let saveIcon;
    let companyIcon;
    let wageIcon;
    let workTimeIcon;
    let addressIcon;

    async function getIcons() {
        saveIcon = await getSvg(save);
        companyIcon = await getSvg(company);
        wageIcon = await getSvg(wage);
        workTimeIcon = await getSvg(workTime);
        addressIcon = await getSvg(workAddress);
    }

    function setBold(id) {
        document.getElementById(id).style.fontWeight = "bold";
        isEditing = true;
    }

    function setRedBorder(id) {
        document.getElementById(id).style.border = "1px solid red";
    }

    let previous;
    let next;

    /*
        async function onPreviousClick() {
            if (previous !== null) {
                let resp = await getAllAds(getCookie("token"), previous, true)
                allAdv = resp[0]
                next = resp[2]
                previous = resp[1]
            }
        }


     */
    async function loadMore() {
        if (next !== null) {
            let resp = await getAllAds(getCookie("token"), next);
            allAdv = allAdv.concat(resp[0]);
            console.log(allAdv);
            next = resp[2];
            previous = resp[1];
        }
    }

    let allAdv = [];
    // Need to use onMount to execute after the DOM is ready
    onMount(async () => {
        await fillUserInfo();
        let resp = await getAllAds(getCookie("token"));
        allAdv = resp[0];
        next = resp[2];
        previous = resp[1];
        await getIcons();
        window.onscroll = function (ev) {
            if (
                window.innerHeight + window.scrollY >=
                document.body.offsetHeight
            ) {
                if (next !== "Null") {
                    loadMore();
                }
            }
        };
    });
</script>

<div style="display: contents">
    <Layout class="h-full">
        <Layout.Header class="static z-0">
            <h1 class="m-3">JobBoard</h1>
            <Layout.Header.Extra slot="extra">
                {#if userConnected}
                    <Button on:click={editProfile} type="ghost">Profile</Button>
                    <Button on:click={logout} class="m-2" type="primary">
                        Logout
                    </Button>
                {:else}
                    <Button
                        on:click={() => showModal("modalLogIn")}
                        class="m-2"
                        type="primary"
                    >
                        Log in
                    </Button>
                    <Button
                        on:click={() => showModal("modalSignIn")}
                        class="m-2"
                        type="primary"
                    >
                        Sign in
                    </Button>
                {/if}
            </Layout.Header.Extra>
        </Layout.Header>
        <Layout.Content>
            <Layout.Content.Body>
                {#each allAdv as advertisement}
                    <div class="md:px-40 py-5 sm:p-5 px-2">
                        <Card class="h-fit">
                            <Card.Cover class="h-32" slot="cover">
                                <img
                                    src={advertisement.companyLogoURL}
                                    alt="cover"
                                    class="object-scale-down h-full w-full"
                                />
                            </Card.Cover>
                            <Card.Content slot="content">
                                <Media>
                                    <Media.Content class="w-full">
                                        <Media.Content.Title>
                                            <h2 class="break-words">
                                                {advertisement.title}
                                            </h2>
                                            <div
                                                class="grid sm:grid-cols-1 md:grid-cols-4"
                                            >
                                                <div class="col col-span-2">
                                                    <div class="flex flex-row">
                                                        <Icon
                                                            data={companyIcon}
                                                            class="mt-3 overflow-visible h-7 w-7"
                                                        />
                                                        <p class="ml-3 m-4">
                                                            {advertisement.companyName}
                                                        </p>
                                                    </div>
                                                </div>
                                                <div class="col col-span-1">
                                                    <div class="flex flex-row">
                                                        <Icon
                                                            data={wageIcon}
                                                            class="mt-3 overflow-visible h-7 w-7"
                                                        />
                                                        <p class="ml-3 m-4">
                                                            {advertisement.wage}
                                                        </p>
                                                    </div>
                                                </div>
                                                <div class="col col-span-1">
                                                    <div class="flex flex-row">
                                                        <Icon
                                                            data={workTimeIcon}
                                                            class="mt-3 overflow-visible h-7 w-7"
                                                        />
                                                        <p
                                                            class="ml-3 m-4 break-normal"
                                                        >
                                                            {advertisement.workTime}
                                                            h
                                                        </p>
                                                    </div>
                                                </div>
                                            </div>
                                        </Media.Content.Title>
                                        <Media.Content.Description>
                                            <Accordion>
                                                <Accordion.Item
                                                    open={open ===
                                                        advertisement.id}
                                                >
                                                    <Accordion.Item.Title
                                                        slot="title"
                                                        on:click={() =>
                                                            handleClick(
                                                                advertisement.id
                                                            )}
                                                    >
                                                        Learn more
                                                    </Accordion.Item.Title>
                                                    <Accordion.Item.Content
                                                        slot="content"
                                                        class="p-5"
                                                    >
                                                        {advertisement.description}
                                                        <div class="grid">
                                                            <div
                                                                class="flex flex-row"
                                                            >
                                                                <Icon
                                                                    data={addressIcon}
                                                                    class="mt-3"
                                                                />
                                                                <p
                                                                    class="ml-1 m-3"
                                                                >
                                                                    {advertisement.address},
                                                                    <span
                                                                        class="break-normal"
                                                                        >{advertisement.city}</span
                                                                    >
                                                                    ({advertisement.zipCode})
                                                                </p>
                                                            </div>
                                                            <p>
                                                                Siren : {advertisement.companySiren}
                                                            </p>
                                                            <Button
                                                                disabled={advertisement.applied}
                                                                type="primary"
                                                                {loading}
                                                                on:click={() =>
                                                                    openModalApply(
                                                                        advertisement.id
                                                                    )}
                                                            >
                                                                {#if advertisement.applied}
                                                                    Applied
                                                                {:else}
                                                                    Apply
                                                                {/if}
                                                            </Button>
                                                        </div>
                                                    </Accordion.Item.Content>
                                                </Accordion.Item>
                                            </Accordion>
                                        </Media.Content.Description>
                                    </Media.Content>
                                </Media>
                            </Card.Content>
                        </Card>
                    </div>
                {/each}
            </Layout.Content.Body>
        </Layout.Content>
    </Layout>
    {#if allModalStates["modalApply"]}
        <Portal>
            <Modal handleClose={closeModalApply}>
                {#if !userConnected}
                    {(allModalStates.modalLogIn = true)}
                {/if}
                <Modal.Content slot="content">
                    <Modal.Content.Header slot="header"
                        >Application form</Modal.Content.Header
                    >
                    {#if noMessage}
                        <Alert type="warn">
                            <Alert.Title slot="title"
                                >Please enter a message</Alert.Title
                            >
                        </Alert>
                    {/if}
                    <Modal.Content.Body slot="body">
                        <Input bind:value={message} name="message">
                            <Input.Label slot="label"
                                >Enter your message for this application</Input.Label
                            >
                        </Input>
                        <Button
                            type="primary"
                            on:click={() => {
                                if (message === "") {
                                    noMessage = true;
                                } else {
                                    apply(idAvert, message);
                                    closeModalApply();
                                }
                            }}
                            >Apply
                        </Button>
                        <Button type="ghost" on:click={closeModalApply}
                            >Cancel</Button
                        >
                        {#if next === "Null"}
                            <h3>No more result</h3>
                        {/if}
                    </Modal.Content.Body>
                </Modal.Content>
            </Modal>
        </Portal>
    {/if}
    {#if allModalStates["modalLogIn"]}
        <Portal>
            <Modal handleClose={hideAllModals}>
                <Modal.Content slot="content">
                    <Modal.Content.Header slot="header"
                        >Login to your account</Modal.Content.Header
                    >
                    {#if badTry}
                        <Alert type="error">
                            <Alert.Title slot="title"
                                >Bad credentials</Alert.Title
                            >
                        </Alert>
                    {/if}
                    <Modal.Content.Body slot="body">
                        <Input bind:value={login} name="login">
                            <Input.Label slot="label">email</Input.Label>
                        </Input>
                        <Input
                            bind:value={password}
                            type="password"
                            name="password"
                        >
                            <Input.Label slot="label">password</Input.Label>
                        </Input>
                        <Button type="primary" on:click={tryLogin}>Go</Button>
                        <Button
                            class="m-2"
                            type="primary"
                            on:click={openModalSignIn}
                            >Don't have account ? Create one
                        </Button>
                        <Button
                            class="m-2"
                            type="ghost"
                            on:click={hideAllModals}>Cancel</Button
                        >
                    </Modal.Content.Body>
                </Modal.Content>
            </Modal>
        </Portal>
    {/if}
    {#if allModalStates["modalSignIn"]}
        <Portal>
            <Modal handleClose={hideAllModals}>
                <Modal.Content slot="content">
                    <Modal.Content.Header slot="header"
                        >Create a new account</Modal.Content.Header
                    >
                    {#if creationFailed}
                        <Alert type="error">
                            <Alert.Title slot="title"
                                >Creation failed please retry</Alert.Title
                            >
                        </Alert>
                    {/if}
                    <Modal.Content.Body slot="body">
                        <form>
                            <Input
                                bind:value={newEmail}
                                id="email"
                                name="login"
                                on:change={() => {
                                    if (
                                        !/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(
                                            newEmail
                                        )
                                    ) {
                                        setRedBorder("email");
                                    } else {
                                        document.getElementById(
                                            "email"
                                        ).style.border = "1px solid green";
                                    }
                                }}
                            >
                                <Input.Label slot="label">email</Input.Label>
                            </Input>
                            <Input
                                bind:value={newPassword}
                                id="password"
                                type="password"
                                name="password"
                                on:input={() => {
                                    if (newPassword.length < 6) {
                                        setRedBorder("password");
                                    } else {
                                        document.getElementById(
                                            "password"
                                        ).style.border = "1px solid green";
                                    }
                                }}
                            >
                                <Input.Label slot="label">password</Input.Label>
                            </Input>
                            <Input bind:value={newName} type="text" name="name">
                                <Input.Label slot="label">name</Input.Label>
                            </Input>
                            <Input
                                bind:value={newSurname}
                                type="text"
                                name="surname"
                            >
                                <Input.Label slot="label">surname</Input.Label>
                            </Input>
                            <Input
                                bind:value={newTel}
                                type="text"
                                id="tel"
                                name="tel"
                                on:change={() => {
                                    if (
                                        !/^[+](\d{3})\)?(\d{3})(\d{5,6})$|^(\d{10,10})$/.test(
                                            newTel
                                        )
                                    ) {
                                        setRedBorder("tel");
                                    } else {
                                        document.getElementById(
                                            "tel"
                                        ).style.border = "1px solid green";
                                    }
                                }}
                            >
                                <Input.Label slot="label">Tel</Input.Label>
                            </Input>
                            <DatePicker
                                max={Date.now()}
                                bind:value={newBirthDate}
                                name="date"
                            >
                                <DatePicker.Label slot="label"
                                    >birthdate</DatePicker.Label
                                >
                            </DatePicker>
                            <Button
                                class="m-2"
                                type="primary"
                                htmlType="submit"
                                on:click={tryCreateUser}
                                >Go
                            </Button>
                            <Button
                                class="m-2"
                                type="ghost"
                                on:click={hideAllModals}>Cancel</Button
                            >
                        </form>
                    </Modal.Content.Body>
                </Modal.Content>
            </Modal>
        </Portal>
    {/if}

    {#if allModalStates["modalEdit"]}
        <Portal>
            <Modal handleClose={closeModalProfile}>
                <Modal.Content slot="content">
                    <Modal.Content.Header slot="header"
                        >Your profile</Modal.Content.Header
                    >
                    {#if editError}
                        <Alert type="error">
                            <Alert.Title slot="title"
                                >Bad credentials</Alert.Title
                            >
                        </Alert>
                    {/if}
                    {#if errorPwdInvalid}
                        <Alert type="error">
                            <Alert.Title slot="title"
                                >Your password is invalid</Alert.Title
                            >
                        </Alert>
                    {/if}
                    <Modal.Content.Body slot="body">
                        <div class="grid sm:grid-cols-1 md:grid-cols-2">
                            <Input
                                id="name"
                                class="col m-3"
                                bind:value={userName}
                                on:input={() => setBold("name")}
                                name="name"
                            >
                                <Input.Label slot="label">name</Input.Label>
                            </Input>
                            <Input
                                id="surname"
                                on:input={() => setBold("surname")}
                                class="col m-3"
                                bind:value={userSurname}
                                name="surname"
                            >
                                <Input.Label slot="label">surname</Input.Label>
                            </Input>
                            <DatePicker
                                id="date"
                                on:input={() => setBold("date")}
                                class="col m-3"
                                bind:value={userBirthDate}
                                name="date"
                            >
                                <DatePicker.Label slot="label"
                                    >birthdate</DatePicker.Label
                                >
                            </DatePicker>
                            <Input
                                id="email"
                                on:input={() => setBold("email")}
                                class="col m-3"
                                bind:value={userEmail}
                                name="email"
                            >
                                <Input.Label slot="label">email</Input.Label>
                            </Input>
                            <Input
                                id="tel"
                                on:input={() => setBold("tel")}
                                class="col m-3"
                                bind:value={userTel}
                                name="tel"
                            >
                                <Input.Label slot="label">telephone</Input.Label
                                >
                            </Input>
                            <Input
                                class="col m-3"
                                type="password"
                                bind:value={userPassword}
                                name="pwd"
                            >
                                <Input.Label slot="label">password</Input.Label>
                            </Input>
                            <Button
                                class="m-2"
                                type="primary"
                                disabled={!isEditing}
                                on:click={() => {
                                    tryEdit();
                                }}
                            >
                                <Icon data={saveIcon} />
                            </Button>
                            <Button
                                class="mt-1"
                                type="ghost"
                                on:click={closeModalProfile}
                            >
                                Cancel
                            </Button>
                        </div>
                    </Modal.Content.Body>
                </Modal.Content>
            </Modal>
        </Portal>
    {/if}
</div>
