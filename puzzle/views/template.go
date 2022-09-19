package views

var PuzzleTemplate = `<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.title}}</title>
    <link href="{{.favicon}}" rel="icon"/>
	<meta name="csrf-token" content={{index .csrfToken}}>
    <link rel="stylesheet" href="/assets/lambda/fonts/roboto/roboto.css?family=Roboto:400,300,100,100italic,300italic,400italic,700,700italic">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/vendor.css"}}">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/paper.css"}}">
    <link rel="stylesheet" href="/assets/lambda/fonts/flaticons/flaticons.css">
    <link rel="stylesheet" href="/assets/lambda/fonts/themify/themify-icons.css">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/moqup.css"}}">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/dataform.css"}}">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/datagrid.css"}}">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/datasource.css"}}">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/chart.css"}}">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/agent.css"}}">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/krud.css"}}">
    <link rel="stylesheet" href="{{index .mix "/assets/lambda/css/puzzle.css"}}">
</head>
<body>
<noscript>To run this application, JavaScript is required to be enabled.</noscript>

<script src="{{index .mix "/assets/lambda/js/manifest.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/vendor.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/datagrid-vendor.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/paper.js"}}"></script>
<script>
    window.app_logo = null
</script>
<div id="puzzle" class="app-wrapper"></div>

<script type="text/javascript" src="/vendor/echart/echarts-en.js"></script>
<script type="text/javascript" src="/vendor/ckeditor/ckeditor.js"></script>
<script>
    window.lambda = {{.lambda_config}};
    window.init = {
        user: {{.User}},
        app_logo:{{.app_logo}},
        app_text:{{.app_text}},
        dbSchema: {{.dbSchema}},
        gridList: {{.gridList}},
        user_fields: {{.user_fields}},
        user_roles: {{.user_roles}},
        data_form_custom_elements: {{.data_form_custom_elements}},
        data_grid_custom_elements: {{.data_grid_custom_elements}}
    };
</script>
<script src="{{index .mix "/assets/lambda/js/moqup.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/dataform.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/dataform-builder.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/datagrid-vendor.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/datagrid.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/datagrid-builder.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/datasource.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/chart.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/chart-builder.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/krud.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/agent.js"}}"></script>
<script src="{{index .mix "/assets/lambda/js/puzzle.js"}}"></script>
</body>
</html>
`
