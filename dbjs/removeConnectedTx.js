function foo(p) {
    let db = require('internal').db,
        toCol = db._collection(p.outColName),
        edgeCol = db._collection(p.edgeColName),
        doc = toCol.document(p.docKey),
        outEdges = edgeCol.outEdges (doc);

    if (outEdges != null && outEdges.length > 0)
        throw ("Deleting vertices with the presence outgoing edges is not allowed");

    edgeCol.inEdges(doc).forEach((edge) => {
        edgeCol.remove (edge);
    });

    toCol.remove (doc);
    return {success: true, entityKey: doc._key};
}
