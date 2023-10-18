        <script>
            import {Icon, Card,Media,Button,Accordion,Layout,Modal, Portal,Input,Alert,DatePicker   } from 'stwui';
            import {getAllAds, Advertisement, submitApply, sendCredentials, createUser, getCookie, getMyId} from '$lib';
            import '../app.postcss';
            import { onMount } from 'svelte';
            import dayjs from "dayjs";
            let loading = false;
            let userConnected = false;
            let idAvert = -1;
            let idUser = 1;//:TODO get id user
            let message ="";
            let noMessage = false;
            let newEmail = "";
            let newPassword = "";
            let newName = "";
            let newSurname = "";
            let newTel = "";
            let creationFailed = false;
            let newBirthDate = new Date();

            let allModalStates = {
                modalApply :false,
                modalLogIn : false,
                modalSignIn : false,
            }
            function setOnlyOneModal(modalName){
                for (const [key, value] of Object.entries(allModalStates)) {
                    if(key!==modalName){
                        allModalStates[key]=false;
                    }else{
                        allModalStates[key]=true;
                    }
                }
                console.log(allModalStates);
            }
            function desacAllModals(){
                for(let modal in allModalStates){
                    allModalStates[modal]=false;
                }
                loading = false;
            }

            async function apply(idAvert,message){
                let token = getCookie("token")
                await submitApply(message,47,idAvert,token);
            }
            async function tryCreateUser() {
                let OK = await createUser(newEmail, newPassword, newName, newSurname, newTel, dayjs(newBirthDate).format("YYYY-MM-DD"));
                if (OK) {
                    toggleUser();
                    let account = await sendCredentials(newEmail, newPassword)
                    idUser = account.UserID
                    let token = account.AuthToken
                    document.cookie = "token=" + token;
                    desacAllModals();
                    badTry = false;
                } else {
                    creationFailed = true;
                }
            }
            let badTry = false;
            function toggleUser() {
                userConnected = !userConnected;
            }
            function toggleLoading() {
                loading = !loading;
            }
            let login = "";
            let password = "";
            async function tryLogin(){
                let account = await sendCredentials(login,password)
                if (account){
                    toggleUser();
                    idUser = account.UserID
                    let token = account.AuthToken
                    document.cookie = "token="+token;
                    desacAllModals();
                    badTry=false;
                }
                else{
                    badTry=true;
                }
            }

            function openModalApply(id) {
                toggleLoading();
                setOnlyOneModal("modalApply");
                idAvert = id;
            }
            function openModalSignIn() {
                setOnlyOneModal("modalSignIn");
            }
            function closeModalApply() {
                toggleLoading();
                desacAllModals();
            }

            let open = '';

            function handleClick(item) {
                if (open === item) {
                    open= '';
                } else {
                    open = item;
                }
            }
            </script>

        <div style="display: contents">
            <Layout class="h-full">
                <Layout.Header class="static z-0">
                    Header
                    <Layout.Header.Extra slot="extra">
                        {#if userConnected}
                        <Button on:click={toggleUser} type="primary">
                                Logout
                        </Button>
                        {:else}
                        <Button on:click={()=>(setOnlyOneModal("modalLogIn"))} type="primary">
                            Log in
                        </Button>
                        <Button on:click={()=>(setOnlyOneModal("modalSignIn"))} type="primary">
                            Sign in
                        </Button>
                        {/if}
                    </Layout.Header.Extra>
                </Layout.Header>
                <Layout.Content>
                <Layout.Content.Body class="grid grid-cols-4 gap-5">
                    {#await getAllAds()}
                        <p>loading...</p>
                    {:then ads}
                    {#each ads as advertisement,i}
                    <div class="col-span-1"></div>
                        <div class="col-span-2">
                            <Card class="h-fit">
                                <Card.Cover class="h-32" slot="cover">
                                    <img
                                        src={advertisement.companyLogoURL}
                                        alt="cover"
                                        class="object-cover object-center w-full aspect-1"
                                    />
                                </Card.Cover>
                                <Card.Content slot="content">
                                    <Media>
                                        <Media.Content class="w-full">
                                            <Media.Content.Title>{advertisement.title}
                                            </Media.Content.Title>
                                            <Media.Content.Description>
                                                {advertisement.description.slice(0,40)}...
                                                <Accordion>
                                                    <Accordion.Item open={open === advertisement.id}>
                                                        <Accordion.Item.Title slot="title" on:click={() => handleClick(advertisement.id)} >
                                                            learn More
                                                        </Accordion.Item.Title>
                                                        <Accordion.Item.Content slot="content">
                                                            <ul>
                                                                <li>
                                                                    {advertisement.companyName}
                                                                </li>
                                                                <li>
                                                                    {advertisement.description}
                                                                </li>
                                                                <li>
                                                                    {advertisement.wage + '€/month'}
                                                                </li>
                                                                <li>
                                                                    {advertisement.address} | {advertisement.city} | {advertisement.zipCode}
                                                                </li>
                                                                <li>
                                                                    {advertisement.workTime + 'H/day'}
                                                                </li>
                                                            </ul>
                                                            <Button type="primary" {loading} on:click={()=>openModalApply(advertisement.id)}>Apply </Button>
                                                        </Accordion.Item.Content>
                                                    </Accordion.Item>
                                                </Accordion>
                                            </Media.Content.Description>

                                        </Media.Content>
                                    </Media>
                                </Card.Content>
                            </Card>
                    </div>
                    <div class="col-span-1"></div>
                    {/each}
                    {/await}
                    </Layout.Content.Body>
                </Layout.Content>
            </Layout>
            {#if allModalStates["modalApply"]}
                        <Portal>
                            <Modal handleClose={closeModalApply}>
                                {#if !userConnected}
                                {setOnlyOneModal("modalLogIn")}
                                {/if}
                                {setOnlyOneModal("modalApply")}
                                <Modal.Content slot="content">
                                    <Modal.Content.Header slot="header">Application form</Modal.Content.Header>
                                    {#if noMessage}
                                        <Alert type="warn">
                                            <Alert.Title slot="title">Please enter a message</Alert.Title>
                                        </Alert>
                                    {/if}
                                    <Modal.Content.Body slot="body">
                                        <Input bind:value={message} name="message">
                                            <Input.Label slot="label">Enter your message for this application</Input.Label>
                                        </Input>
                                        <Button type="primary" on:click={()=>{
                                                if(message == ""){
                                                    noMessage = true;
                                                }else{
                                                    apply(idAvert,message);
                                                    closeModalApply();
                                                }
                                            }
                                        }>Apply</Button>
                                            <Button type="ghost" on:click={closeModalApply}>Cancel</Button>
                                    </Modal.Content.Body>
                                </Modal.Content>
                            </Modal>
                        </Portal>
            {/if}
            {#if allModalStates["modalLogIn"]}

                        <Portal>
                            <Modal handleClose={desacAllModals}>
                                <Modal.Content slot="content">
                                    <Modal.Content.Header slot="header">Login to your account</Modal.Content.Header>
                                    {#if badTry}
                                        <Alert type="error">
                                            <Alert.Title slot="title">Bad credentials</Alert.Title>
                                        </Alert>
                                    {/if}
                                    <Modal.Content.Body slot="body">
                                            <Input bind:value={login} name="login">
                                                <Input.Label slot="label">email</Input.Label>
                                            </Input>
                                            <Input bind:value={password} type="password" name="password">
                                                <Input.Label slot="label">password</Input.Label>
                                            </Input>
                                            <Button type="primary" on:click={tryLogin}>Go
                                            </Button>
                                            <Button type="primary" on:click={openModalSignIn}>Don't have account ? create one
                                            </Button>
                                            <Button type="ghost" on:click={desacAllModals}>Cancel</Button>

                                    </Modal.Content.Body>
                                </Modal.Content>
                            </Modal>
                        </Portal>
            {/if}
            {#if allModalStates["modalSignIn"]}
            <Portal>
                <Modal handleClose={desacAllModals}>
                    <Modal.Content slot="content">
                        <Modal.Content.Header slot="header">Create a new account</Modal.Content.Header>
                        {#if creationFailed}
                            <Alert type="error">
                                <Alert.Title slot="title">TODO</Alert.Title>
                            </Alert>
                        {/if}
                        <Modal.Content.Body slot="body">
                            <form>
                                <Input bind:value={newEmail} name="login">
                                    <Input.Label slot="label">email</Input.Label>
                                </Input>
                                <Input bind:value={newPassword} type="password" name="password">
                                    <Input.Label slot="label">password</Input.Label>
                                </Input>
                                <Input bind:value={newName} type="text" name="name">
                                    <Input.Label slot="label">Name</Input.Label>
                                </Input>
                                <Input bind:value={newSurname} type="text" name="surname">
                                    <Input.Label slot="label">Surname</Input.Label>
                                </Input>
                                <Input bind:value={newTel} type="text" name="tel">
                                    <Input.Label slot="label">Tel</Input.Label>
                                </Input>
                                <DatePicker bind:value={newBirthDate} name="date">
                                    <DatePicker.Label slot="label">Birthdate</DatePicker.Label>
                                </DatePicker>
                                <Button type="primary" on:click={tryCreateUser}>Go
                                </Button>
                                <Button type="ghost" on:click={desacAllModals}>Cancel</Button>
                            </form>
                        </Modal.Content.Body>
                    </Modal.Content>
                </Modal>
            </Portal>
            {/if}
        </div>