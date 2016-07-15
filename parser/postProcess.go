package parser

import "reflect"

type loadExternalUse interface {
	loadExternalUse(conf PostProcessConfig) (err error)
}

var loadExternalUseRef = reflect.TypeOf((*loadExternalUse)(nil)).Elem()

func loadExternalUseExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(loadExternalUse).loadExternalUse(conf)
}

type fixRequiredBySyntax interface {
	fixRequiredBySyntax() (err error)
}

var fixRequiredBySyntaxRef = reflect.TypeOf((*fixRequiredBySyntax)(nil)).Elem()

func fixRequiredBySyntaxExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(fixRequiredBySyntax).fixRequiredBySyntax()
}

type fixDefaultMediaType interface {
	fixDefaultMediaType(conf PostProcessConfig) (err error)
}

var fixDefaultMediaTypeRef = reflect.TypeOf((*fixDefaultMediaType)(nil)).Elem()

func fixDefaultMediaTypeExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(fixDefaultMediaType).fixDefaultMediaType(conf)
}

type fixEmptyAnnotation interface {
	fixEmptyAnnotation() (err error)
}

var fixEmptyAnnotationRef = reflect.TypeOf((*fixEmptyAnnotation)(nil)).Elem()

func fixEmptyAnnotationExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(fixEmptyAnnotation).fixEmptyAnnotation()
}

type fixAnnotationBracket interface {
	fixAnnotationBracket() (err error)
}

var fixAnnotationBracketRef = reflect.TypeOf((*fixAnnotationBracket)(nil)).Elem()

func fixAnnotationBracketExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(fixAnnotationBracket).fixAnnotationBracket()
}

type fillProperties interface {
	fillProperties(library Library) (err error)
}

var fillPropertiesRef = reflect.TypeOf((*fillProperties)(nil)).Elem()

func fillPropertiesExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(fillProperties).fillProperties(conf.Library())
}

type fillTrait interface {
	fillTrait(library Library) (err error)
}

var fillTraitRef = reflect.TypeOf((*fillTrait)(nil)).Elem()

func fillTraitExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(fillTrait).fillTrait(conf.Library())
}

type fillURIParams interface {
	fillURIParams() (err error)
}

var fillURIParamsRef = reflect.TypeOf((*fillURIParams)(nil)).Elem()

func fillURIParamsExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(fillURIParams).fillURIParams()
}

type fillExample interface {
	fillExample(conf PostProcessConfig) (err error)
}

var fillExampleRef = reflect.TypeOf((*fillExample)(nil)).Elem()

func fillExampleExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(fillExample).fillExample(conf)
}

type checkTypoError interface {
	checkTypoError() (err error)
}

var checkTypoErrorRef = reflect.TypeOf((*checkTypoError)(nil)).Elem()

func checkTypoErrorExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(checkTypoError).checkTypoError()
}

type checkUnusedAnnotation interface {
	checkUnusedAnnotation(annotationUsage map[string]bool) (err error)
}

var checkUnusedAnnotationRef = reflect.TypeOf((*checkUnusedAnnotation)(nil)).Elem()

func checkUnusedAnnotationExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(checkUnusedAnnotation).checkUnusedAnnotation(conf.AnnotationUsage())
}

type afterCheckUnusedAnnotation interface {
	afterCheckUnusedAnnotation(conf PostProcessConfig) (err error)
}

var afterCheckUnusedAnnotationRef = reflect.TypeOf((*afterCheckUnusedAnnotation)(nil)).Elem()

func afterCheckUnusedAnnotationExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(afterCheckUnusedAnnotation).afterCheckUnusedAnnotation(conf)
}

type checkUnusedTrait interface {
	checkUnusedTrait(traitUsage map[string]bool) (err error)
}

var checkUnusedTraitRef = reflect.TypeOf((*checkUnusedTrait)(nil)).Elem()

func checkUnusedTraitExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(checkUnusedTrait).checkUnusedTrait(conf.TraitUsage())
}

type afterCheckUnusedTrait interface {
	afterCheckUnusedTrait(conf PostProcessConfig) (err error)
}

var afterCheckUnusedTraitRef = reflect.TypeOf((*afterCheckUnusedTrait)(nil)).Elem()

func afterCheckUnusedTraitExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(afterCheckUnusedTrait).afterCheckUnusedTrait(conf)
}

type checkExample interface {
	checkExample(conf PostProcessConfig) (err error)
}

var checkExampleRef = reflect.TypeOf((*checkExample)(nil)).Elem()

func checkExampleExec(v interface{}, conf PostProcessConfig) (err error) {
	return v.(checkExample).checkExample(conf)
}

type postProcessFunc func(v interface{}, conf PostProcessConfig) (err error)

var postProcessInfoMap = map[reflect.Type]postProcessFunc{
	loadExternalUseRef:            loadExternalUseExec,
	fixRequiredBySyntaxRef:        fixRequiredBySyntaxExec,
	fixDefaultMediaTypeRef:        fixDefaultMediaTypeExec,
	fixEmptyAnnotationRef:         fixEmptyAnnotationExec,
	fixAnnotationBracketRef:       fixAnnotationBracketExec,
	fillPropertiesRef:             fillPropertiesExec,
	fillTraitRef:                  fillTraitExec,
	fillURIParamsRef:              fillURIParamsExec,
	fillExampleRef:                fillExampleExec,
	checkTypoErrorRef:             checkTypoErrorExec,
	checkUnusedAnnotationRef:      checkUnusedAnnotationExec,
	afterCheckUnusedAnnotationRef: afterCheckUnusedAnnotationExec,
	checkUnusedTraitRef:           checkUnusedTraitExec,
	afterCheckUnusedTraitRef:      afterCheckUnusedTraitExec,
	checkExampleRef:               checkExampleExec,
}

func postProcess(v interface{}, conf PostProcessConfig) (err error) {
	implements := []reflect.Type{
		loadExternalUseRef,
		fixRequiredBySyntaxRef,
		fixDefaultMediaTypeRef,
		fixEmptyAnnotationRef,
		fixAnnotationBracketRef,
		fillPropertiesRef,
		fillTraitRef,
		fillURIParamsRef,
		fillExampleRef,
		checkTypoErrorRef,
		checkUnusedAnnotationRef,
		afterCheckUnusedAnnotationRef,
		checkUnusedTraitRef,
		afterCheckUnusedTraitRef,
		checkExampleRef,
	}
	for _, implement := range implements {
		if err = postProcessImplement(reflect.ValueOf(v), implement, conf); err != nil {
			return
		}
	}
	return
}

var reflectTypeValue = reflect.TypeOf(Value{})
var reflectTypeValuePtr = reflect.TypeOf(&Value{})
var reflectTypeLibrary = reflect.TypeOf(Library{})
var reflectTypeLibraryPtr = reflect.TypeOf(&Library{})

func postProcessImplement(val reflect.Value, implement reflect.Type, conf PostProcessConfig) (err error) {
	switch val.Type() {
	case reflectTypeValue, reflectTypeValuePtr:
		// no need to post process Value
		return nil
	case reflectTypeLibrary:
		conf = newPostProcessConfig(conf.RootDocument(), val.Interface().(Library), conf.Parser())
	case reflectTypeLibraryPtr:
		conf = newPostProcessConfig(conf.RootDocument(), *val.Interface().(*Library), conf.Parser())
	}

	if v := queryPostProcessImplement(val, implement); v != nil {
		if err = postProcessInfoMap[implement](v, conf); err != nil {
			return
		}
	}

	kind := val.Kind()
	if kind == reflect.Ptr {
		if val.IsNil() {
			return
		}
		kind = val.Elem().Kind()
		val = val.Elem()
	}

	switch kind {
	case reflect.Struct:
		for i, n := 0, val.NumField(); i < n; i++ {
			if err = postProcessImplement(val.Field(i), implement, conf); err != nil {
				return
			}
		}
	case reflect.Slice:
		for i, n := 0, val.Len(); i < n; i++ {
			if err = postProcessImplement(val.Index(i), implement, conf); err != nil {
				return
			}
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			if err = postProcessImplement(val.MapIndex(key), implement, conf); err != nil {
				return
			}
		}
	}

	return
}

// queryPostProcessImplement return not nil if val can run implement
func queryPostProcessImplement(val reflect.Value, implement reflect.Type) interface{} {
	if val.CanAddr() {
		addr := val.Addr()
		if addr.CanInterface() && addr.Type().Implements(implement) {
			return addr.Interface()
		}
	}
	if val.CanInterface() && val.Type().Implements(implement) {
		return val.Interface()
	}
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		elem := val.Elem()
		if elem.CanInterface() && elem.Type().Implements(implement) {
			return val.Elem().Interface()
		}
	}
	return nil
}
