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
            type: "integer",
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
            type: "integer",
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

function createUsCommentDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "usCommentListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/userstory/{0}/comment"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/userstory/{0}/addcomment"},
        ],
        transformRequest: function (dsRequest) {
            switch (dsRequest.operationType) {
                case "add":
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
                    url = operationBinding.dataURL.format(dsRequest.originalData.usKey);
                    break;

                default:
                    url = operationBinding.dataURL.format(dsRequest.originalData._key);
            }

            return url;
        },
        fields: [{
            name: "title",
            validators: []
        }, {
            name: "text",
            validators: []
        }, {
            name: "createDate",
            type: 'datetime',
            validators: []
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }, {
            name: "user.displayName"
        }, {
            name: "usKey",
            foreignKey: "usListDS._key",
            hidden: true
        }]
    });
}

function createUSCommentGrid () {
    createUsCommentDS();

    return isc.ListGrid.create({
        ID: "usCommentList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: usCommentListDS,
        autoFetchData: false,
        canEdit: false,
        wrapCells: true,
        cellHeight: 60,
        fields:[{
            name: "title",
            title: "Заголовок",
            width: 100
        }, {
            name: "text",
            title: "Текст",
            wrap: true,
            width: 300
        }, {
            name: "createDate",
            title: "Создан",
            width: 100,
            type: 'datetime'
        }, {
            name: "user.displayName",
            title: "Создал",
            width: 200
        }],
        editComplete: function() {
            refreshRelatedGrid(usList.getSelectedRecord(), usList, usCommentList);
        }
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
        titleOrientation: "top",
        show: function (values) {
            textEditForm.setValues(values);
            textEditForm.setTitleOrientation("top");
            this.Super("show", arguments)
        },
        items: [isc.DynamicForm.create ({
            ID: "textEditForm",
            fields: [
                {name:"text", title: "Текст", type:"textArea", width: 300, height: 200},
                {
                    title:"OK",
                    type:"button",
                    click: function () {
                        var values = textEditForm.getValues()
                        var url = "/api/v.0.0.1/userstory/" + usList.getSelectedRecord()._key + "/addtext";
                        isc.RPCManager.sendRequest({ data: JSON.stringify(values), callback: function (data) {
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

    return isc.ListGrid.create({
        ID: "usList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: usListDS,
        autoFetchData: false,
        canEdit: true,
        canExpandRecords: true,
        expansionMode: "detailField",
        detailField: "text",
        selectionUpdated : function (data) {
            refreshRelatedGrid(usList.getSelectedRecord(), usList, usCommentList);
            refreshRelatedGrid(usList.getSelectedRecord(), usList, tsList);
        },
        editComplete: function() {
            refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
        }
    });
}

function createFeatureGrid () {
    createFeatureDS ();

    isc.ListGrid.create({
        ID: "featureList",
        height: "25%",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: featureListDS,
        autoFetchData: false,
        canEdit: true,
        title: "Фичи",
        selectionUpdated : function (data) {
            refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
        },
        editComplete: function() {
            refreshRelatedGrid(processList.getSelectedRecord(), processList, featureList);
        }
    });

    var splitPane = isc.SectionStack.create ({
        width: "30%",
        showResizeBar: true,
        visibilityMode: "multiple",
        sections: [{
            expanded: true,
            title: "Feature",
            items: [
                featureList
            ],
            controls: [
                isc.ImgButton.create({
                    src: "[SKIN]/actions/add.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Добавить",
                    click: function() {
                        if (!processList.getSelectedRecord())
                            return false;

                        featureList.startEditingNew({processKey: processList.getSelectedRecord()._key});
                        return false;
                    }
                }),
                isc.ImgButton.create({
                    autoDraw: false,
                    src: "[SKIN]actions/remove.png",
                    size: 16,
                    showFocused: false, showRollOver: false, showDown: false,
                    prompt: "Удалить",
                    click: function () {
                        featureList.removeSelectedData(function () {
                            refreshRelatedGrid(processList.getSelectedRecord(), processList, featureList);
                        })
                    }
                })
            ]
        }, {
            expanded: true,
            title: "User Story",
            items: [
                createUsGrid()
            ],
            controls: [
                isc.ImgButton.create({
                    src: "[SKIN]/actions/add.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Добавить стори",
                    click: function() {
                        if (!featureList.getSelectedRecord())
                            return false;

                        usList.startEditingNew({featureKey: featureList.getSelectedRecord()._key, serial: usList.getData().totalRows + 1});
                        return false;
                    }
                }),
                isc.ImgButton.create({
                    src: "[SKIN]/actions/remove.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Удалить стори",
                    click: function () {
                        usList.removeSelectedData(function () {
                            refreshRelatedGrid(featureList.getSelectedRecord(), featureList, usList);
                        })
                    }
                }),
                isc.ImgButton.create({
                    src: "[SKIN]/actions/edit.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Добавить текст к сторе",
                    click: function () {
                        if (!usList.getSelectedRecord())
                            return false;

                        textEditWindow.show({text: usList.getSelectedRecord().text});
                    }
                })
            ]
        }, {
            expanded: true,
            title: "Комментарии",
            items: [
                createUSCommentGrid()
            ],
            controls: [
                isc.ImgButton.create({
                    src: "[SKIN]/actions/add.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Добавить комментарий",
                    click: function() {
                        if (!usList.getSelectedRecord())
                            return false;

                        usCommentList.startEditingNew({usKey: usList.getSelectedRecord()._key});
                        return false;
                    }
                })
            ]
        }]
    });

    return splitPane;
}
