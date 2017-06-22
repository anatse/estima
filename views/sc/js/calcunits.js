function createCalcUnitsSelectWindow () {
    createCalcUnitDS ();
    createTsCalcUnitDS ();

    isc.Window.create({
        ID: "cuSelectWindow",
        title: "Привязка функциональных точек",
        width: 1024,
        height: 400,
        autoCenter: true,
        isModal: true,
        showModalMask: true,
        autoDraw: false,
        titleOrientation: "center",
        show: function (values) {
            calcUnitList.fetchData();
            if (!tsList.getSelectedRecord())
                return false;

            tsCalcUnitList.fetchData({"tsKey": tsList.getSelectedRecord()._key})
            tsCalcUnitList.refreshData();
            this.Super("show", arguments)
        },
        items: [
            isc.VLayout.create ({
                members: [
                    isc.ListGrid.create({
                        height: "50%",
                        ID: "calcUnitList",
                        alternateRecordStyles:true,
                        showAllRecords:true,
                        dataSource: calcUnitListDS,
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
                                title: "Привязать ф.точку",
                                click: function () {
                                    if (calcUnitList.getSelection().length == 0) {
                                        isc.warn ("Не выбрано ни одной ф.точки");
                                        return false;
                                    }

                                    var sel = calcUnitList.getSelection();
                                    for (var i=0;i<sel.length;i++) {
                                        tsCalcUnitList.addData ({
                                            tsKey: tsList.getSelection()[0]._key,
                                            _key: sel[i]._key
                                        });
                                    }

                                    tsCalcUnitList.refreshData();
                                }
                            }),
                            isc.Button.create ({
                                showRollOver: false,
                                showFocus: false,
                                showDown: false,
                                showSelected: false,
                                width: 150,
                                icon: "[SKIN]NavigationBar/miniNav~2_up.png",
                                title: "Отвязать ф.точку",
                                click: function () {
                                    if (tsCalcUnitList.getSelection().length == 0) {
                                        isc.warn ("Не выбрано ни одной ф.точки для отвязки");
                                        return false;
                                    }

                                    var sel = tsCalcUnitList.getSelection();
                                    for (var i=0;i<sel.length;i++) {
                                        var data = sel[i]
                                        console.log(data);
                                        tsCalcUnitList.removeData (data);
                                    }

                                    tsCalcUnitList.refreshData();
                                }
                            })
                        ]
                    }),
                    isc.ListGrid.create({
                        height: "100%",
                        ID: "tsCalcUnitList",
                        alternateRecordStyles:true,
                        showAllRecords:true,
                        dataSource: tsCalcUnitListDS,
                        autoFetchData: false,
                        selectionAppearance:"checkbox",
                        canEdit: true,
                        editComplete: function() {
                            tsCalcUnitList.refreshData();
                        }
                    })
                ]
            })
        ]
    });
}

function createCalcUnitDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "calcUnitListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/cu/list"}
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
            name: "version",
            width: 200,
            title: "Версия"
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

function createTsCalcUnitDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "tsCalcUnitListDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/ts/{0}/listCu"},
            {operationType: "add", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/ts/{0}/addCu"},
            {operationType: "remove", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/ts/{0}/cu/{1}/remove"},
            {operationType: "update", dataProtocol: "postMessage", requestProperties: {httpMethod: "POST"}, dataURL: "/api/v.0.0.1/ts/{0}/cu/{1}/update"}
        ],
        transformRequest: function (dsRequest) {
            switch (dsRequest.operationType) {
                case "add":
                case "update":
                case "remove":
                    dsRequest.data.complexity = parseFloat(dsRequest.data.complexity);
                    dsRequest.data.newFlag = parseFloat(dsRequest.data.newFlag);
                    return JSON.stringify(dsRequest.data);

                default:
                    return dsRequest.data;
            }
        },
        getDataURL: function (dsRequest) {
            var operationBinding = this.getOperationBinding(dsRequest);
            var url = "";
            if (!dsRequest.originalData.tsKey)
                dsRequest.originalData.tsKey = tsList.getSelectedRecord()._key;

            switch (dsRequest.operationType) {
                case "fetch":
                case "add":
                    url = operationBinding.dataURL.format(dsRequest.originalData.tsKey);
                    break;

                default:
                    url = operationBinding.dataURL.format(dsRequest.originalData.tsKey, dsRequest.originalData._key);
            }

            return url;
        },
        fields: [{
            name: "name",
            title: "Название",
            width: 200,
            canEdit: false,
            validators: []
        }, {
            name: "description",
            width: 300,
            title: "Описание",
            canEdit: false,
            validators: []
        }, {
            name: "weight",
            type: "integer",
            width: 60,
            title: "Кол-во"
        }, {
            name: "newFlag",
            valueMap: {
                1: "Новая",
                0.5: "Модификация"
            },
            width: 100,
            title: "Переисп"
        }, {
            name: "complexity",
            valueMap: {
                0.5: "Простая",
                1: "Типовая",
                1.2: "Сложная"
            },
            width: 100,
            title: "Сложность"
        }, {
            name: "extCoef",
            type: "double",
            width: 70,
            title: "Доп. коэффициент"
        }, {
            name: "version",
            type: "text",
            canEdit: false,
            title: "Версия"
        }, {
            name: "status",
            type: "integer",
            canEdit: false,
            title: "Статус"
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }, {
            name: "tsKey",
            foreignKey: "tsListDS._key",
            hidden: true
        }]
    });
}

createCalcUnitsSelectWindow ();