function foo(params) {
    let db = require('internal').db,
        toCol = db._collection (params.toColName),
        edgeCol = db._collection(params.edgeColName),
        doc = db._collection(params.fromColName).document(params.fKey),
        activeText = null,
        numActives = 0;

    db._createStatement({
        query: "FOR v, e IN OUTBOUND '" + doc._id + "' " + params.edgeColName + " FILTER e.label == 'text' && v.active RETURN v"
    }).execute().toArray().forEach((d) => {
        activeText = d;
        numActives++;
    });

    if (numActives > 1)
        throw ("Found more than one active texts for given feature. Please, contact with administrator to fix problem");

    // Set inactive for previous text
    activeText.active = false;
    toCol.save (activeText);

    // Save new text with ++version
    params.text.active = true;
    params.text.version = activeText.version + 1
    let toDoc = toCol.save (params.text);

    // Add edge between text and object
    edgesCol.save ({_from: doc._id, _to: toDoc._id, label: 'text'});
    return {success: true, entityKey: toDoc._key};
}
