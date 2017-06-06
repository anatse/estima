function createFeatureDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "featureListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/process/{0}/feature/list"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/process/{0}/feature/add"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/feature/{0}/update"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/feature/{0}/remove"},
        ],
        transformRequest: function (dsRequest) {
            switch (dsRequest.operationType) {
                case "add":
                case "update":
                case "remove":
                    return JSON.stringify(dsRequest.data, function (key, value) {
                        return value;
                    });
                    break;

                default:
                    return dsRequest.data;
            }
        },
        getDataURL: function (dsRequest) {
            var operationBinding = this.getOperationBinding(dsRequest);
            var url = "";
            switch (dsRequest.operationType) {
                case "add":
                case "fetch":
                    url = operationBinding.dataURL.format(dsRequest.originalData.processKey);
                    break;

                default:
                    url = operationBinding.dataURL.format(dsRequest.originalData._key);
            }

            return url;
        },
        fields: [{
            name: "name",
            title: "Название",
            width: 200,
            validators: []
        }, {
            name: "description",
            width: 300,
            title: "Описание",
            validators: []
        }, {
            name: "status",
            title: "Статус"
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }, {
            name: "processKey",
            foreignKey: "stageProcessListDS._key",
            hidden: true
        }]
    });
}

function createUSDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "usListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/feature/{0}/userstory/list"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/feature/{0}/userstory/add"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/userstory/{0}/update"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/userstory/{0}/remove"},
        ],
        transformRequest: function (dsRequest) {
            switch (dsRequest.operationType) {
                case "add":
                case "update":
                case "remove":
                    return JSON.stringify(dsRequest.data, function (key, value) {
                        return value;
                    });
                    break;

                default:
                    return dsRequest.data;
            }
        },
        getDataURL: function (dsRequest) {
            var operationBinding = this.getOperationBinding(dsRequest);
            var url = "";
            switch (dsRequest.operationType) {
                case "add":
                case "fetch":
                    url = operationBinding.dataURL.format(dsRequest.originalData.featureKey);
                    break;

                default:
                    url = operationBinding.dataURL.format(dsRequest.originalData._key);
            }

            return url;
        },
        fields: [{
            name: "name",
            title: "Название",
            width: 200,
            validators: []
        }, {
            name: "who",
            width: 300,
            title: "Кто",
            validators: []
        },{
            name: "what",
            width: 300,
            title: "Что",
            validators: []
        },{
            name: "why",
            width: 300,
            title: "Зачем",
            validators: []
        },{
            name: "description",
            width: 300,
            title: "Описание",
            validators: []
        }, {
            name: "status",
            title: "Статус"
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }, {
            name: "serial",
            title: "Номер",
            type: "integer"
        }, {
            name: "text",
            hidden: true
        }, {
            name: "featureKey",
            foreignKey: "featureListDS._key",
            hidden: true
        }]
    });
}

function createWindowText () {
    isc.Window.create({
        ID: "textEditWindow",
        title: "Text Edit Window",
        autoSize: true,
        autoCenter: true,
        isModal: true,
        showModalMask: true,
        autoDraw: false,
        show: function (values) {
            userEditForm.setValues(values);
            this.Super("show", arguments)
        },
        items: [isc.DynamicForm.create ({
            ID: "textEditForm",
            fields: [
                {name:"text", type:"textArea"},
                {
                    title:"OK",
                    type:"button",
                    click: function () {
                        var values = textEditForm.getValues()
                        var url = "/api/v.0.0.1/userstory/" + usList.getSelectedRecord()._key + "/addtext";
                        isc.RPCManager.sendRequest({ data: JSON.stringify(values), callback: function (data) {
                            console.log(arguments)
                            refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
                            textEditWindow.close();
                        }, actionURL: url, httpMethod: 'POST', contentType: "application/json",
                            useSimpleHttp: true});
                    }
                }
            ],
        })]
    });
}

function createUsGrid () {
    createUSDS();
    createWindowText ();

    var editControls = isc.ToolStrip.create({
        members: [
            isc.LayoutSpacer.create({ width:"*" }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/add.png",
                prompt: "Add record",
                click: function() {
                    if (!featureList.getSelectedRecord())
                        return false;

                    usList.startEditingNew({featureKey: featureList.getSelectedRecord()._key, serial: usList.getData().totalRows + 1});
                    return false;
                }
            }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/remove.png",
                prompt: "Remove selected record",
                click: function () {
                    usList.removeSelectedData(function () {
                        refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
                    })
                }
            }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/edit.png",
                prompt: "Add text to the selected record",
                click: function () {
                    if (!usList.getSelectedRecord())
                        return false;

                    textEditWindow.show();
                }
            })
        ]
    });

    return isc.ListGrid.create({
        ID: "usList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: usListDS,
        autoFetchData: false,
        canEdit: true,
        gridComponents:[editControls, "header", "body"],
        canExpandRecords: true,
        expansionMode: "detailField",
        detailField: "text",
        selectionUpdated : function (data) {
            console.log (data);
        },
        editComplete: function() {
            refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
        }
    });
}

function createFeatureGrid () {
    createFeatureDS ();

    var editControls = isc.ToolStrip.create({
        members: [
            isc.LayoutSpacer.create({ width:"*" }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/add.png",
                prompt: "Add record",
                click: function() {
                    if (!processList.getSelectedRecord())
                        return false;

                    featureList.startEditingNew({processKey: processList.getSelectedRecord()._key});
                    return false;
                }
            }),
            isc.ToolStripButton.create({
                icon: "[SKIN]/actions/remove.png",
                prompt: "Remove selected record",
                click: function () {
                    featureList.removeSelectedData(function () {
                        refreshRelatedGrid(processList.getSelectedRecord(), processList, featureList);
                    })
                }
            })
        ]
    });

    var tabSet = isc.TabSet.create({
        autoDraw: false,
        tabBarPosition: "top",
        tabs: [
            {
                title: "User Stories",
                pane: createUsGrid()
            },
            {
                title: "Feature components"
            }
        ]
    });

    isc.ListGrid.create({
        ID: "featureList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: featureListDS,
        autoFetchData: false,
        canEdit: true,
        gridComponents:[editControls, "header", "body"],
        selectionUpdated : function (data) {
            refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
        },
        editComplete: function() {
            refreshRelatedGrid(processList.getSelectedRecord(), processList, featureList);
        }
    });

    return isc.HLayout.create({
        // width: "100%",
        // height: "100%",
        members: [
            featureList,
            tabSet
        ]
    });
}
