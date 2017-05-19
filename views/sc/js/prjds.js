'use strict';

isc.Log.setPriority("Log", isc.Log.DEBUG);
isc.Log.logDebug ('test messge');

isc.DataSource.create({
    dataURL:"/user/projects",
    dataFormat:"json",
    allowAdvancedCriteria: true,
    ID:"userProjectListDS",
    childrenField: "body",
    recordXpath: "//body",
    fields:[{
            name:"number",
            valueXPath:"body/number",
            title: "Номер проекта"
        }, {
            name:"description",
            valueXPath:"body/description",
            title:"Описание",
            validators:[
            ]
        }, {
            name:"name",
            title:"Название",
            valueXPath:"body/name",
            validators:[
            ]
        }, {
            name:"status",
            valueXPath:"body/stratus",
            title: "Статус"
        }, {
            name:"startDate",
            valueXPath:"body/startDate",
            type: "date",
            format: "",
            title: "Дата старта"
        }, {
            name:"endDate",
            valueXPath:"body/endDate",
            title: "Дата окончания"
        }, {
            name: "flag",
            valueXPath:"body/flag",
            title: "флаг"
        }
    ]
});

isc.ListGrid.create({
    ID: "userProjectList",
    width:500, height:224, alternateRecordStyles:true, showAllRecords:true,
    dataSource: userProjectListDS,
    autoFetchData: true
});