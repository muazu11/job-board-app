<script>
    import { Button, Layout, Tabs, Card, Table, Pagination } from "stwui";

    const tabs = [
        {
            title: "Users",
            key: "users",
        },
        {
            title: "Companies",
            key: "companies",
        },
        {
            title: "Advertisements",
            key: "advertisements",
        },
        {
            title: "Applications",
            key: "applications",
        },
    ];

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
                <Card bordered={false} class="h-[calc(100vh-14rem)]">
                    <Card.Header
                        slot="header"
                        class="font-bold text-lg flex justify-between items-center py-3"
                    >
                        Card Header
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
                                {#each results as item}
                                    <Table.Body.Row id={item.id}>
                                        <Table.Body.Row.Cell column={0}
                                            >{item.email}</Table.Body.Row.Cell
                                        >
                                        <Table.Body.Row.Cell column={1}
                                            >{item.name}</Table.Body.Row.Cell
                                        >
                                        <Table.Body.Row.Cell column={2}
                                            >{item.surname}</Table.Body.Row.Cell
                                        >
                                        <Table.Body.Row.Cell column={3}
                                            >{item.date_of_birth}</Table.Body.Row.Cell
                                        >
                                        <Table.Body.Row.Cell column={4}
                                            >{item.role}</Table.Body.Row.Cell
                                        >
                                        <Table.Body.Row.Cell column={5}
                                            >{item.password}</Table.Body.Row.Cell
                                        >
                                    </Table.Body.Row>
                                {/each}
                            </Table.Body>
                            <!--  <Table.Footer slot="footer">
                                <Pagination
                                    start={data.start}
                                    end={data.end}
                                    total={data.total}
                                    currentPage={parseInt(currentPage)}
                                    {onPreviousClick}
                                    {onNextClick}
                                    {onPageClick}
                                />
                            </Table.Footer>-->
                        </Table>
                    </Card.Content>
                </Card>
            </Layout.Content.Body>
        </Layout.Content>
    </Layout>
</div>
