function foo(params) {
    let db = require('internal').db,
        toCol = db._collection (params.toColName),
        edgeCol = db._collection(params.edgeColName),
        doc = db._collection(params.fromColName).document(params.fKey),
        activeText = null,
        numActives = 0;

    // Save new text wit ++version
    let toDoc = toCol.save (params.comment);

    // Add edge between comment and object
    edgesCol.save ({_from: doc._id, _to: toDoc._id, label: 'comment'});

    // Add edge between comment and user
    edgesCol.save ({_from: params.userId, _to: toDoc._id, label: 'comment'});

    return {success: true, entityKey: toDoc._key};
}
