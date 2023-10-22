<script>
    import { onMount } from "svelte";
    import {
        getCookie,
        getAllUsers,
        getAllCompanies,
        getAllAdvertisements,
        getAllApplications,
    } from "$lib";
    import { Button, Layout, Card, Table } from "stwui";

    let userColumns = [
        { label: "email", placement: "left" },
        { label: "name", placement: "left" },
        { label: "surname", placement: "left" },
        { label: "date_of_birth", placement: "left" },
        { label: "role", placement: "left" },
        { label: "password", placement: "left" },
    ];

    let companyColumns = [
        { label: "name", placement: "left" },
        { label: "siren", placement: "left" },
        { label: "logo_url", placement: "left" },
    ];

    let advertisementColumns = [
        { label: "company_id", placement: "left" },
        { label: "title", placement: "left" },
        { label: "description", placement: "left" },
        { label: "wage", placement: "left" },
        { label: "address", placement: "left" },
        { label: "zip_code", placement: "left" },
        { label: "city", placement: "left" },
        { label: "work_time", placement: "left" },
    ];

    let applicationColumns = [
        { label: "advertisement_id", placement: "left" },
        { label: "applicant_id", placement: "left" },
        { label: "message", placement: "left" },
        { label: "created_at", placement: "left" },
    ];

    let results = [
        {
            id: 1,
            email: "test@test.com",
            name: "test",
            surname: "testos",
            date_of_birth: "21-03-2001",
            role: "admin",
            password: "......",
        },
        {
            id: 2,
            email: "testi@test.com",
            name: "testi",
            surname: "testos",
            date_of_birth: "22-03-2001",
            role: "user",
            password: "......",
        },
    ];

    let token;
    let users = [];
    let companies = [];
    let advertisements = [];
    let applications = [];

    onMount(async () => {
        token = getCookie("token");

        let userResp = await getAllUsers(token);
        users = userResp["data"];

        let companyResp = await getAllCompanies(token);
        companies = companyResp["data"];

        let advertisementResp = await getAllAdvertisements(token);
        advertisements = advertisementResp["data"];

        let applicationResp = await getAllApplications(token);
        applications = applicationResp["data"];
        console.log(applications);
    });
</script>

<div style="display: contents">
    <Layout class="h-full">
        <Layout.Header class="static z-0">
            <h1 class="m-3">JobBoard</h1>
            <Layout.Header.Extra slot="extra">
                <Button href="../" class="m-2" type="primary">Home</Button>
            </Layout.Header.Extra>
        </Layout.Header>
        <Layout.Content>
            <Layout.Content.Body>
                <Card bordered={false} class="m-8">
                    <Card.Header
                        slot="header"
                        class="font-bold text-lg flex justify-between items-center py-3"
                    >
                        Users
                        <Button slot="extra" type="primary">New Item</Button>
                    </Card.Header>
                    <Card.Content
                        slot="content"
                        class="p-0 sm:p-0"
                        style="height: calc(100% - 64px);"
                    >
                        <Table
                            class="rounded-md overflow-hidden h-full"
                            columns={userColumns}
                        >
                            <Table.Header slot="header" />
                            <Table.Body slot="body">
                                {#each users as user}
                                    <Table.Body.Row id={user.id}>
                                        <Table.Body.Row.Cell column={0}>
                                            {user.email}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={1}>
                                            {user.name}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={2}>
                                            {user.surname}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={3}>
                                            {user.dateOfBirthUTC}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={4}>
                                            {user.role}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={5}>
                                            {user.password}
                                        </Table.Body.Row.Cell>
                                    </Table.Body.Row>
                                {/each}
                            </Table.Body>
                        </Table>
                    </Card.Content>
                </Card>

                <Card bordered={false} class="m-8">
                    <Card.Header
                        slot="header"
                        class="font-bold text-lg flex justify-between items-center py-3"
                    >
                        Companies
                        <Button slot="extra" type="primary">New Item</Button>
                    </Card.Header>
                    <Card.Content
                        slot="content"
                        class="p-0 sm:p-0"
                        style="height: calc(100% - 64px);"
                    >
                        <Table
                            class="rounded-md overflow-hidden h-full"
                            columns={companyColumns}
                        >
                            <Table.Header slot="header" />
                            <Table.Body slot="body">
                                {#each companies as company}
                                    <Table.Body.Row id={company.id}>
                                        <Table.Body.Row.Cell column={0}>
                                            {company.name}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={1}>
                                            {company.siren}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={2}>
                                            {company.logoURL}
                                        </Table.Body.Row.Cell>
                                    </Table.Body.Row>
                                {/each}
                            </Table.Body>
                        </Table>
                    </Card.Content>
                </Card>

                <Card bordered={false} class="m-8">
                    <Card.Header
                        slot="header"
                        class="font-bold text-lg flex justify-between items-center py-3"
                    >
                        Advertisements
                        <Button slot="extra" type="primary">New Item</Button>
                    </Card.Header>
                    <Card.Content
                        slot="content"
                        class="p-0 sm:p-0"
                        style="height: calc(100% - 64px);"
                    >
                        <Table
                            class="rounded-md overflow-hidden h-full"
                            columns={advertisementColumns}
                        >
                            <Table.Header slot="header" />
                            <Table.Body slot="body">
                                {#each advertisements as advertisement}
                                    <Table.Body.Row id={advertisement.id}>
                                        <Table.Body.Row.Cell column={0}>
                                            {advertisement.companyID}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={1}>
                                            {advertisement.title}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={2}>
                                            {advertisement.description}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={3}>
                                            {advertisement.wage}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={4}>
                                            {advertisement.adress}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={5}>
                                            {advertisement.zipCode}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={6}>
                                            {advertisement.city}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={7}>
                                            {advertisement.workTimeNs}
                                        </Table.Body.Row.Cell>
                                    </Table.Body.Row>
                                {/each}
                            </Table.Body>
                        </Table>
                    </Card.Content>
                </Card>

                <Card bordered={false} class="m-8">
                    <Card.Header
                        slot="header"
                        class="font-bold text-lg flex justify-between items-center py-3"
                    >
                        Applications
                        <Button slot="extra" type="primary">New Item</Button>
                    </Card.Header>
                    <Card.Content
                        slot="content"
                        class="p-0 sm:p-0"
                        style="height: calc(100% - 64px);"
                    >
                        <Table
                            class="rounded-md overflow-hidden h-full"
                            columns={applicationColumns}
                        >
                            <Table.Header slot="header" />
                            <Table.Body slot="body">
                                {#each applications as application}
                                    <Table.Body.Row id={application.id}>
                                        <Table.Body.Row.Cell column={0}>
                                            {application.advertisementID}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={1}>
                                            {application.applicantID}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={2}>
                                            {application.message}
                                        </Table.Body.Row.Cell>
                                        <Table.Body.Row.Cell column={3}>
                                            {application.createdAt}
                                        </Table.Body.Row.Cell>
                                    </Table.Body.Row>
                                {/each}
                            </Table.Body>
                        </Table>
                    </Card.Content>
                </Card>
            </Layout.Content.Body>
        </Layout.Content>
    </Layout>
</div>
