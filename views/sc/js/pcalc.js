function createComponentsSelectWindow () {
    createComponentDS ();
    createFeatureComponentDS ();

    isc.Window.create({
        ID: "cmpSelectWindow",
        title: "Выбор компонентов",
        width: 900,
        height: 400,
        autoCenter: true,
        isModal: true,
        showModalMask: true,
        autoDraw: false,
        titleOrientation: "center",
        show: function (values) {
            componentList.fetchData();
            if (!featureList.getSelectedRecord())
                return false;

            featureComponentList.fetchData({"featureKey": featureList.getSelectedRecord()._key})
            this.Super("show", arguments)
        },
        items: [
            isc.VLayout.create ({
                members: [
                    isc.ListGrid.create({
                        height: "50%",
                        ID: "componentList",
                        alternateRecordStyles:true,
                        showAllRecords:true,
                        dataSource: componentListDS,
                        autoFetchData: false,
                        selectionAppearance:"checkbox",
                        canEdit: false
                    }),
                    isc.HLayout.create ({
                        layoutAlign:"center",
                        align: "center",
                        members: [
                            isc.Button.create ({
                                showRollOver: false,
                                showFocus: false,
                                showDown: false,
                                showSelected: false,
                                width: 150,
                                icon: "[SKIN]NavigationBar/miniNav~2_down.png",
                                title: "Привязать компонент",
                                click: function () {
                                    if (componentList.getSelection().length == 0) {
                                        isc.warn ("Не выбрано ни одного компонента");
                                        return false;
                                    }

                                    var sel = componentList.getSelection();
                                    for (var i=0;i<sel.length;i++) {
                                        featureComponentList.addData ({
                                            featureKey: featureList.getSelection()[0]._key,
                                            _key: sel[i]._key
                                        });
                                    }

                                    featureComponentList.refreshData();
                                }
                            }),
                            isc.Button.create ({
                                showRollOver: false,
                                showFocus: false,
                                showDown: false,
                                showSelected: false,
                                width: 150,
                                icon: "[SKIN]NavigationBar/miniNav~2_up.png",
                                title: "Отвязать компонент",
                                click: function () {
                                    if (featureComponentList.getSelection().length == 0) {
                                        isc.warn ("Не выбрано ни одного компонента для отвязки");
                                        return false;
                                    }

                                    var sel = featureComponentList.getSelection();
                                    for (var i=0;i<sel.length;i++) {
                                        var data = sel[i]
                                        console.log(data);
                                        featureComponentList.removeData (data);
                                    }

                                    featureComponentList.refreshData();
                                }
                            })
                        ]
                    }),
                    isc.ListGrid.create({
                        height: "100%",
                        ID: "featureComponentList",
                        alternateRecordStyles:true,
                        showAllRecords:true,
                        dataSource: featureComponentListDS,
                        autoFetchData: false,
                        selectionAppearance:"checkbox",
                        canEdit: false
                    })
                ]
            })
        ]
    });
}

function createComponentDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "componentListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/component/list"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/component/add"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/component/{id}/update"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/component/{id}/remove"},
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
                case "update":
                case "remove":
                    url = operationBinding.dataURL.format(dsRequest.originalData._key);
                    break;

                default:
                    url = operationBinding.dataURL;
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
            name: "owner",
            width: 200,
            title: "Владелец"
        }, {
            name: "dueDate",
            type: "datetime",
            title: "Дата выхода"
        }, {
            name: "status",
            type: "integer",
            title: "Статус"
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }]
    });
}

function createFeatureComponentDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "featureComponentListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/feature/{0}/listComponent"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/feature/{0}/addComponent"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/feature/{0}/component/{1}/remove"},
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
                case "fetch":
                case "add":
                    url = operationBinding.dataURL.format(dsRequest.originalData.featureKey);
                    break;

                default:
                    url = operationBinding.dataURL.format(dsRequest.originalData.featureKey, dsRequest.originalData._key);
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
            name: "owner",
            width: 200,
            title: "Владелец"
        }, {
            name: "dueDate",
            type: "datetime",
            title: "Дата выхода"
        }, {
            name: "status",
            type: "integer",
            title: "Статус"
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }, {
            name: "featureKey",
            foreignKey: "featureListDS._key",
            hidden: true
        }]
    });
}

createComponentsSelectWindow ();