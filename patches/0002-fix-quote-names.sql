diff --git a/pkg/sql/sem/tree/type_name.go b/pkg/sql/sem/tree/type_name.go
index 898009a..5d4423a 100644
--- a/pkg/sql/sem/tree/type_name.go
+++ b/pkg/sql/sem/tree/type_name.go
@@ -56,8 +56,7 @@ func (t *TypeName) String() string {
 
 // SQLString implements the ResolvableTypeReference interface.
 func (t *TypeName) SQLString() string {
-	// FmtBareIdentifiers prevents the TypeName string from being wrapped in quotations.
-	return AsStringWithFlags(t, FmtBareIdentifiers)
+	return AsStringWithFlags(t, FmtSimple)
 }
 
 // FQString renders the type name in full, not omitting the prefix
@@ -250,14 +249,12 @@ func (node *ArrayTypeReference) Format(ctx *FmtCtx) {
 
 // SQLString implements the ResolvableTypeReference interface.
 func (node *ArrayTypeReference) SQLString() string {
-	// FmtBareIdentifiers prevents the TypeName string from being wrapped in quotations.
-	return AsStringWithFlags(node, FmtBareIdentifiers)
+	return AsStringWithFlags(node, FmtSimple)
 }
 
 // SQLString implements the ResolvableTypeReference interface.
 func (name *UnresolvedObjectName) SQLString() string {
-	// FmtBareIdentifiers prevents the TypeName string from being wrapped in quotations.
-	return AsStringWithFlags(name, FmtBareIdentifiers)
+	return AsStringWithFlags(name, FmtSimple)
 }
 
 // IsReferenceSerialType returns whether the input reference is a known
