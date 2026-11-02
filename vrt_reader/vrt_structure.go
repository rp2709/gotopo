package vrt_reader

import "encoding/xml"

type VRTDataset struct {
	XMLName      xml.Name        `xml:"VRTDataset"`
	RasterXSize  int             `xml:"rasterXSize,attr"`
	RasterYSize  int             `xml:"rasterYSize,attr"`
	SRS          SRS             `xml:"SRS"`
	GeoTransform string          `xml:"GeoTransform"`
	RasterBands  []VRTRasterBand `xml:"VRTRasterBand"`
}

type SRS struct {
	DataAxisToSRSAxisMapping string `xml:"dataAxisToSRSAxisMapping,attr"`
	Value                    string `xml:",chardata"`
}

type VRTRasterBand struct {
	DataType    string          `xml:"dataType,attr"`
	Band        int             `xml:"band,attr"`
	Metadata    Metadata        `xml:"Metadata"`
	NoDataValue string          `xml:"NoDataValue"`
	ColorInterp string          `xml:"ColorInterp"`
	Sources     []ComplexSource `xml:"ComplexSource"`
}

type Metadata struct {
	Items []MetadataItem `xml:"MDI"`
}

type MetadataItem struct {
	Key   string `xml:"key,attr"`
	Value string `xml:",chardata"`
}

type ComplexSource struct {
	SourceFilename SourceFilename `xml:"SourceFilename"`
	SourceBand     int            `xml:"SourceBand"`
	SrcRect        Rect           `xml:"SrcRect"`
	DstRect        Rect           `xml:"DstRect"`
	NoData         string         `xml:"NODATA"`
}

type SourceFilename struct {
	RelativeToVRT bool   `xml:"relativeToVRT,attr"`
	Value         string `xml:",chardata"`
}

type Rect struct {
	XOff  float64 `xml:"xOff,attr"`
	YOff  float64 `xml:"yOff,attr"`
	XSize float64 `xml:"xSize,attr"`
	YSize float64 `xml:"ySize,attr"`
}
