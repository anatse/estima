function createTSDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "tsListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/userstory/{0}/tstory/list"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/userstory/{0}/tstory/add"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/tstory/{0}/update"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/tstory/{0}/remove"},
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
                    url = operationBinding.dataURL.format(dsRequest.originalData.usKey);
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
            name: "text",
            hidden: true
        }, {
            name: "usKey",
            foreignKey: "usListDS._key",
            hidden: true
        }]
    });
}

function createTsCommentDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "tsCommentListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/tstory/{0}/comment"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/tstory/{0}/addcomment"},
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
                    url = operationBinding.dataURL.format(dsRequest.originalData.tsKey);
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
            name: "tsKey",
            foreignKey: "tsListDS._key",
            hidden: true
        }]
    });
}

function createTSCommentGrid () {
    createTsCommentDS();

    return isc.ListGrid.create({
        ID: "tsCommentList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: tsCommentListDS,
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
            refreshRelatedGrid(tsList.getSelectedRecord(), tsList, tsCommentList);
        }
    });
}

function createTsWindowText () {
    isc.Window.create({
        ID: "tsTextEditWindow",
        title: "Text Edit Window",
        autoSize: true,
        autoCenter: true,
        isModal: true,
        showModalMask: true,
        autoDraw: false,
        titleOrientation: "top",
        show: function (values) {
            tsTextEditForm.setValues(values);
            tsTextEditForm.setTitleOrientation("top");
            this.Super("show", arguments)
        },
        items: [isc.DynamicForm.create ({
            ID: "tsTextEditForm",
            fields: [
                {name:"text", title: "Текст", type:"textArea", width: 300, height: 200},
                {
                    title:"OK",
                    type:"button",
                    click: function () {
                        var values = tsTextEditForm.getValues()
                        var url = "/api/v.0.0.1/tstory/" + tsList.getSelectedRecord()._key + "/addtext";
                        isc.RPCManager.sendRequest({ data: JSON.stringify(values), callback: function (data) {
                            refreshRelatedGrid(usList.getSelectedRecord(), usList, tsList);
                            tsTextEditWindow.close();
                        }, actionURL: url, httpMethod: 'POST', contentType: "application/json",
                            useSimpleHttp: true});
                    }
                }
            ],
        })]
    });
}

function createTsGrid () {
    createTSDS();
    createTsWindowText ();

    return isc.ListGrid.create({
        ID: "tsList",
        showResizeBar: true,
        alternateRecordStyles:true,
        showAllRecords:true,
        dataSource: tsListDS,
        autoFetchData: false,
        canEdit: true,
        canExpandRecords: true,
        expansionMode: "detailField",
        detailField: "text",
        selectionUpdated : function (data) {
            refreshRelatedGrid(tsList.getSelectedRecord(), tsList, tsCommentList);
        },
        editComplete: function() {
            refreshRelatedGrid(usList.getSelectedRecord(), usList, tsList);
        }
    });
}

function createTslayput () {
    return isc.SectionStack.create ({
        width: "30%",
        showResizeBar: true,
        visibilityMode: "multiple",
        sections: [{
            expanded: true,
            title: "Техистории",
            items: [
                createTsGrid()
            ],
            controls: [
                isc.ImgButton.create({
                    src: "[SKIN]/actions/add.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Добавить техсторю",
                    click: function() {
                        if (!usList.getSelectedRecord())
                            return false;

                        tsList.startEditingNew({usKey: usList.getSelectedRecord()._key});
                        return false;
                    }
                }),
                isc.ImgButton.create({
                    src: "[SKIN]/actions/remove.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Удалить  техстори",
                    click: function () {
                        tsList.removeSelectedData(function () {
                            refreshRelatedGrid(usList.getSelectedRecord(), usList, tsList);
                        })
                    }
                }),
                isc.ImgButton.create({
                    src: "[SKIN]/actions/edit.png",
                    size: 16,
                    showFocused: false,
                    showRollOver: false,
                    showDown: false,
                    prompt: "Добавить текст к техсторе",
                    click: function () {
                        if (!tsList.getSelectedRecord())
                            return false;

                        tsTextEditWindow.show({text: tsList.getSelectedRecord().text});
                    }
                })
            ]
        }, {
            expanded: true,
            title: "Комментарии",
            items: [
                createTSCommentGrid()
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
                        if (!tsList.getSelectedRecord())
                            return false;

                        tsCommentList.startEditingNew({tsKey: tsList.getSelectedRecord()._key});
                        return false;
                    }
                })
            ]
        }]
    });
}