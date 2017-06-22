function createProjectPriceWindow () {
    createProjectPriceDS ();

    isc.Window.create({
        ID: "projectPriceWindow",
        title: "Расчсет стоимости проекта",
        width: 1400,
        height: 700,
        autoCenter: true,
        isModal: true,
        showModalMask: true,
        autoDraw: false,
        titleOrientation: "center",
        show: function (values) {
            projectPrice.fetchData();
            projectPrice.refreshData();

            this.Super("show", arguments)
        },
        items: [
            isc.VLayout.create ({
                members: [
                    isc.ListGrid.create({
                        ID: "projectPrice",
                        alternateRecordStyles:true,
                        showAllRecords:true,
                        dataSource: projectPriceDS,
                        autoFetchData: false,
                        canEdit: false,
                        canMultiGroup: true,
                        groupByField: [
                            "prjName",
                            // "stageName",
                            // "prcName",
                            // "featureName",
                            // "ustoryName",
                            // "tstoryName",
                            "group"
                        ],
                        groupStartOpen: "all",
                        showGridSummary: true,
                        showGroupSummary: true
                    })
                ]
            })
        ]
    });
}

function createProjectPriceDS () {
    isc.DataSource.create({
        dataFormat: "json",
        allowAdvancedCriteria: true,
        ID: "projectPriceDS",
        recordXPath: "body",
        operationBindings: [
            {operationType: "fetch", dataProtocol: "", requestProperties: {httpMethod: "GET"}, dataURL: "/api/v.0.0.1/calcproject/{0}"}
        ],
        transformRequest: function (dsRequest) {
            switch (dsRequest.operationType) {
                default:
                    return dsRequest.data;
            }
        },
        getDataURL: function (dsRequest) {
            var operationBinding = this.getOperationBinding(dsRequest);
            var url = operationBinding.dataURL.format(userProjectList.getSelectedRecord()._key);
            return url;
        },
        fields: [{
            name: "prjName",
            title: "Проект",
            width: 150,
        }, {
            name: "stageName",
            width: 150,
            title: "Стадия",
        }, {
            name: "prcName",
            width: 150,
            title: "Процесс"
        }, {
            name: "featureName",
            width: 150,
            title: "Фича"
        }, {
            name: "ustoryName",
            width: 150,
            title: "Польз-я история"
        }, {
            name: "ustoryInfo",
            width: 150,
            title: "История п-я"
        }, {
            name: "tstoryName",
            width: 150,
            title: "Техистория"
        }, {
            name: "group",
            width: 150,
            title: "Группа",
            valueMap: {
                "Analyze": "Анализ",
                "Integration": "Интеграция",
                "Development": "Разработка",
                "DevOps": "Поддержка",
                "Architect": "Архитектура",
                "SysArch": "Прикладная архитектура",
                "IntArch": "Интеграционная архитектура"
            }

        }, {
            name: "price",
            width: 150,
            type: "double",
            title: "Стоимость (SP)"
        }, {
            name: "ustoryKey",
            hidden: true
        }, {
            name: "tstoryKey",
            hidden: true
        }, {
            name: "_key",
            primaryKey: true,
            hidden: true
        }]
    });
}

createProjectPriceWindow ();