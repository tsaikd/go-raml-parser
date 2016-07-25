package parser

func mergeAPIType(dst *APIType, fromList ...APIType) {
	for _, from := range fromList {
		mergeTypeDeclaration(&dst.TypeDeclaration, from.TypeDeclaration)

		if dst.ObjectType.IsEmpty() {
			dst.ObjectType = from.ObjectType
		}
		if dst.ScalarType.IsEmpty() {
			dst.ScalarType = from.ScalarType
		}
		if dst.String.IsEmpty() {
			dst.String = from.String
		}
		if dst.IsArray && dst.ArrayType.IsEmpty() {
			dst.ArrayType = from.ArrayType
		}
		if dst.FileType.IsEmpty() {
			dst.FileType = from.FileType
		}

		dst.NativeType = from.NativeType
	}
}

func mergeTypeDeclaration(dst *TypeDeclaration, fromList ...TypeDeclaration) {
	for _, from := range fromList {
		mergeUnimplement(&dst.Default, from.Default)
		mergeUnimplement(&dst.Schema, from.Schema)
		// do not merge Type field because Type should not be empty
		// do not merge Example(s) field because Example will be filled by fillExample()
		if dst.DisplayName == "" {
			dst.DisplayName = from.DisplayName
		}
		if dst.Description == "" {
			dst.Description = from.Description
		}
		dst.Annotations = mergeAnnotations(dst.Annotations, from.Annotations)
		mergeUnimplement(&dst.Facets, from.Facets)
		mergeUnimplement(&dst.XML, from.XML)
	}
}

func mergeAnnotations(dst Annotations, fromList ...Annotations) Annotations {
	if dst == nil {
		dst = Annotations{}
	}
	for _, from := range fromList {
		for name, annotation := range from {
			if dstanno := dst[name]; dstanno == nil {
				dst[name] = annotation
			}
		}
	}
	return dst
}

func mergeUnimplement(dst *Unimplement, fromList ...Unimplement) {
	for _, from := range fromList {
		if dst.IsEmpty() {
			*dst = from
		}
	}
}
