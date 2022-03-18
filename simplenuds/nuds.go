package simplenuds

import (
	"encoding/xml"
	"time"
)

// See http://nomisma.org/nuds.xsd
// NUDS: root element of a NUDS document.
type NUDS struct {
	// See https://pkg.go.dev/encoding/xml for
	// the XML struct tags supported by Golang.

	XMLName xml.Name `xml:"nuds"`

	// See https://github.com/golang/go/issues/11496
	// for discussion of namespace serialization causing us to put namespace like this.
	XMLNS          string `xml:"xmlns,attr"`
	METS_NS        string `xml:"xmlns:mets,attr"`
	TEI_NS         string `xml:"xmlns:tei,attr"`
	XS_NS          string `xml:"xmlns:xs,attr"`
	XLINK_NS       string `xml:"xmlns:xlink,attr"`
	XSI_NS         string `xml:"xmlns:xsi,attr"`
	SchemaLocation string `xml:"xsi:schemaLocation,attr"`

	// The @recordType is a required attribute for the <nuds> root element.
	// A record must be 'conceptual' (a coin type or die typology) or
	// 'physical'.
	RecordType string `xml:"recordType,attr"`

	Control  Control  `xml:"control"`
	DescMeta DescMeta `xml:"descMeta"`
	DigRep   *DigRep  `xml:"digRep"`
}

// The area of the instance that contains control information about its identity,
// creation, maintenance, status, and the rules and authorities used in the composition
// of the description....
type Control struct {
	RecordID string `xml:"recordId"`

	// TODO
	// <xs:element maxOccurs="unbounded" minOccurs="0" ref="otherRecordId"/>

	// <xs:element ref="publicationStatus"/>
	PublicationStatus PublicationStatus `xml:"publicationStatus"`

	// <xs:element ref="maintenanceStatus"/>
	MaintenanceStatus MaintenanceStatus `xml:"maintenanceStatus"`

	// <xs:element ref="maintenanceAgency"/>
	MaintenanceAgency MaintenanceAgency `xml:"maintenanceAgency"`

	// <xs:element ref="maintenanceHistory"/>
	MaintenanceHistory MaintenanceHistory `xml:"maintenanceHistory"`

	// <xs:element ref="rightsStmt"/>
	RightsStmt RightsStmt `xml:"rightsStmt"`

	// <xs:element maxOccurs="unbounded" minOccurs="0" ref="semanticDeclaration"/>
	// <xs:attribute ref="xml:id"/>
	// <xs:attribute ref="xml:lang"/>
}

// The Descriptive Metadata element is one of two required elements within <nuds>.
// It is the container for all descriptive metadata containers for an object or typology.
// <typeDesc> is the only required child element.
type DescMeta struct {
	// <xs:element maxOccurs="unbounded" ref="title"/>
	// Title may be repeated if the title is to be available in numerous languages
	Title []Title `xml:"title"`

	// TODO
	//<xs:element minOccurs="0" ref="descriptionSet"/>
	//<xs:element minOccurs="0" ref="noteSet"/>
	//<xs:element minOccurs="0" ref="subjectSet"/>

	// <xs:element minOccurs="1" ref="typeDesc"/>
	TypeDesc TypeDesc `xml:"typeDesc"`

	//<xs:element minOccurs="0" ref="physDesc"/>
	PhysDesc *PhysDesc `xml:"physDesc"`

	// TODO
	//<xs:element minOccurs="0" ref="undertypeDesc"/>
	//<xs:element minOccurs="0" ref="findspotDesc"/>
	//<xs:element minOccurs="0" ref="refDesc"/>
	//<xs:element minOccurs="0" ref="adminDesc"/>
}

// Object title. It may be repeated if the title is to be available in numerous languages
type Title struct {
	// <xs:attribute ref="xml:id"/>

	// <xs:attribute ref="xml:lang" use="required"/>
	Lang string `xml:"xml:lang,attr"`

	// <xs:attributeGroup ref="nuds_a.localType"/>
	Value string `xml:",chardata"`
}

// The Typological Description, <typeDesc>, is a container for
// typological characteristics of a resource, whether a physical
// object or a coin type. The <typeDesc> is the only required top-level
// descriptive element within <descMeta>.
type TypeDesc struct {
	// <xs:element minOccurs="0" ref="objectType"/>
	// <xs:choice>
	// <xs:element minOccurs="0" ref="date"/>
	// <xs:element minOccurs="0" ref="dateRange"/>
	// </xs:choice>
	// <xs:element minOccurs="0" maxOccurs="1" ref="dateOnObject"/>

	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="denomination"/>
	Denomination []Denomination `xml:"denomination"`

	// <xs:element minOccurs="0" ref="manufacture"/>

	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="material"/>
	Material []Material `xml:"material"`

	// <xs:element minOccurs="0" maxOccurs="1" ref="shape"/>
	// <xs:element minOccurs="0" ref="authority"/>
	// <xs:element minOccurs="0" ref="geographic"/>
	// <xs:element minOccurs="0" ref="obverse"/>
	// <xs:element minOccurs="0" ref="reverse"/>
	// <xs:element minOccurs="0" ref="edge"/>
	// <xs:element minOccurs="0" ref="weightStandard"/>
	// <xs:element minOccurs="0" ref="typeSeries"/>
}

// The Physical Description element of <descMeta> is a container for the physical characteristics
// of an object. It should not be used for typological records.
type PhysDesc struct {
	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="authenticity"/>
	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="originalIntendedUse"/>
	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="peculiarityOfProduction"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="axis"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="channelOrientation"/>
	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="chemicalAnalysis"/>
	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="color"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="conservationState"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="countermark"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="dateOnObject"/>

	// <xs:element minOccurs="0" maxOccurs="1" ref="measurementsSet"/>
	MeasurementsSet *MeasurementsSet `xml:"measurementsSet"`

	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="serialNumber"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="shape"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="testmark"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="watermark"/>
}

// The <measurementsSet> is a container for physical measurments of an object.
type MeasurementsSet struct {
	// <xs:element minOccurs="0" ref="diameter"/>
	Diameter *Diameter `xml:"diameter"`

	// TODO
	// <xs:element minOccurs="0" ref="height"/>
	// <xs:element minOccurs="0" ref="thickness"/>
	// <xs:element minOccurs="0" ref="length"/>
	// <xs:element minOccurs="0" ref="specificGravity"/>
	// <xs:element minOccurs="0" ref="width"/>

	// <xs:element minOccurs="0" ref="weight"/>
	Weight *Weight `xml:"weight"`
}

// Diameter (in decimal numbers) of a round object. Units and precision
// may be included as attributes. The <diameter> should not be used in
// conjunction with <height> and <width>.
type Diameter struct {
	// TODO Try to use Golang anonymous type inheritance here

	// Names the unit used for the measurement. Suggested values include: 1] g; 2] cm; 3] mm
	Units string `xml:"units,attr,omitempty"`

	Value string `xml:",chardata"`
}

// Weight encoded in decimal numbers, usually by grams.
type Weight struct {
	// Example:
	// <weight units="g">4.14</weight>

	Units string `xml:"units,attr,omitempty"`

	Value string `xml:",chardata"`
}

// The <denomination>, usually defined by a Nomisma URI by means of XLink attributes.
// <xs:attributeGroup ref="m.default"/>
// <xs:attributeGroup ref="xlink:simpleLink"/>
type Denomination string

// The <material> (e.g., silver), usually defined by a Nomisma URI by means of XLink attributes.
// For example <material xlink:href="http://nomisma.org/id/ar" xlink:type="simple">Silver</material>
type Material struct {
	// <xs:attributeGroup ref="m.default"/>
	// <xs:attributeGroup ref="xlink:simpleLink"/>
	HRef string `xml:"href,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`

	Text string `xml:",chardata"`
}

type PublicationStatus struct {
	// <xs:enumeration value="inProcess"/>
	// <xs:enumeration value="approved"/>
	// <xs:enumeration value="approvedSubtype"/>
	// <xs:enumeration value="deprecatedType"/>
	Value string `xml:",chardata"`
}

// ... records the current drafting status of an NUDS instance.
type MaintenanceStatus struct {
	// <xs:enumeration value="new"/>
	// <xs:enumeration value="derived"/>
	// <xs:enumeration value="revised"/>
	// <xs:enumeration value="cancelled"/>
	// <xs:enumeration value="cancelledSplit"/>
	// <xs:enumeration value="cancelledReplaced"/>
	Value string `xml:",chardata"`
}

// The institution or service responsible for the creation, maintenance,
// and/or dissemination of the NUDS instance.
type MaintenanceAgency struct {
	// <xs:element ref="agencyName"/>
	AgencyName AgencyName `xml:"agencyName"`

	// TODO
	// <xs:element minOccurs="0" ref="agencyCode"/>
	// <xs:element maxOccurs="unbounded" minOccurs="0" ref="otherAgencyCode"/>
}

// The name of the institution or service responsible for the creation,
// maintenance, and/or dissemination of the NUDS instance.
type AgencyName struct {
	// TODO
	// <xs:attribute ref="xml:id"/>
	// <xs:attribute ref="xml:lang"/>

	Value string `xml:",chardata"`
}

// A required wrapper element within <control> to record the history and
// creation of the EAC-CPF instance.
type MaintenanceHistory struct {
	// <xs:element maxOccurs="unbounded" ref="maintenanceEvent"/>
	MaintenanceEvent []MaintenanceEvent `xml:"maintenanceEvent"`

	// TODO
	// <xs:attribute ref="xml:id"/>
	// <xs:attribute ref="xml:lang"/>
}

// ... details modification dates and descriptions of changes, including the
// person or process responsible for making the change.
type MaintenanceEvent struct {
	// <xs:element ref="eventType"/>
	EventType EventType `xml:"eventType"`

	// <xs:element ref="eventDateTime"/>
	EventDateTime EventDateTime `xml:"eventDateTime"`

	// <xs:element ref="agentType"/>
	AgentType AgentType `xml:"agentType"`

	// <xs:element ref="agent"/>
	Agent Agent `xml:"agent"`

	// TODO
	// <xs:element minOccurs="0" maxOccurs="1" ref="eventDescription"/>
	// <xs:element minOccurs="0" maxOccurs="1" ref="source"/>
	// <xs:attribute ref="xml:id"/>
	// <xs:attribute ref="xml:lang"/>
}

// ... On first creation, the event type would be "created".
// A "derived" event type is available to indicate that the record was
// derived from another descriptive system ...
type EventType struct {
	// <xs:enumeration value="created"/>
	// <xs:enumeration value="revised"/>
	// <xs:enumeration value="deleted"/>
	// <xs:enumeration value="cancelled"/>
	// <xs:enumeration value="cancelledSplit"/>
	// <xs:enumeration value="derived"/>
	// <xs:enumeration value="updated"/>
	Value string `xml:",chardata"`
}

type EventDateTime struct {
	// For example <eventDateTime standardDateTime="2020-07-31T13:12:46-05:00">Fri, 31 Jul 2020</eventDateTime>

	// The date or date and time represented in a standard form for computer processing.
	StandardDateTime string `xml:"standardDateTime,attr"`

	Value string `xml:",chardata"`
}

// The type of agent responsible for modifying the NUDS record: "human" or "machine".
type AgentType struct {
	// <xs:enumeration value="human"/>
	// <xs:enumeration value="machine"/>
	Value string `xml:",chardata"`
}

// The name of the agent responsible for the <maintenanceEvent>.
type Agent struct {
	// e.g. <agent>PHP</agent>
	Value string `xml:",chardata"`
}

// Statement of rights recording the NUDS record
type RightsStmt struct {
	// <xs:element minOccurs="0" ref="copyrightHolder"/>

	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="license"/>
	License []License `xml:"license"`

	// <xs:element minOccurs="0" maxOccurs="unbounded" ref="rights"/>
	// <xs:element minOccurs="0" ref="preferCite"/>
	// <xs:element minOccurs="0" ref="date"/>
}

// A statement of the metadata license. It is recommended to link to a license
// provided by Creative Commons for images or other digital facsimiles or the
// Open Data Commons for the metadata.
type License struct {
	// <xs:attribute name="for">
	For string `xml:"for,attr"`

	// <xs:simpleType>
	Type string `xml:"xlink:type,attr"`

	Href  string `xml:"xlink:href,attr"`
	Value string `xml:",chardata"`
}

type DigRep struct {
	// XMLName  xml.Name `xml:"control"`
	FileSec FileSec `xml:"mets:fileSec"`
}

// See http://www.loc.gov/standards/mets/mets.xsd

// The overall purpose of the content file section element <fileSec> is to
// provide an inventory of and the location for the content files that
// comprise the digital object being described in the METS document.
type FileSec struct {
	FileGrp []FileGrp `xml:"mets:fileGrp"`
}

// A sequence of file group elements <fileGrp> can be used group the digital
// files comprising the content of a METS object either into a flat arrangement or,
// because each file group element can itself contain one or more file group elements,
// into a nested (hierarchical) arrangement....
type FileGrp struct {
	File []File `xml:"mets:file"`
}

// The file element <file> provides access to the content files for the
// digital object being described by the METS document. A <file> element
// may contain one or more <FLocat> elements which provide pointers to a
// content file and/or a <FContent> element which wraps an encoded version
// of the file....
type File struct {
	FLocat []FLocat `xml:"mets:FLocat"`
}

// The file location element <FLocat> provides a pointer to the location
// of a content file. It uses the XLink reference syntax to provide linking
// information indicating the actual location of the content file, along
// with other attributes specifying additional linking information.
type FLocat struct {
	LOCTYPE string `xml:"LOCTYPE,attr"`
	Href    string `xml:"xlink:href,attr"`
}

// Defaulters

func (nuds *NUDS) DefaultDigRep() *DigRep {
	if nuds.DigRep == nil {
		nuds.DigRep = &DigRep{}
	}

	return nuds.DigRep
}

func (descMeta *DescMeta) DefaultPhysDesc() *PhysDesc {
	if descMeta.PhysDesc == nil {
		descMeta.PhysDesc = &PhysDesc{}
	}

	return descMeta.PhysDesc
}

func (physDesc *PhysDesc) DefaultMeasurementsSet() *MeasurementsSet {
	if physDesc.MeasurementsSet == nil {
		physDesc.MeasurementsSet = &MeasurementsSet{}
	}

	return physDesc.MeasurementsSet
}

func (descMeta *DescMeta) DefaultTitle() []Title {
	if descMeta.Title == nil {
		descMeta.Title = []Title{
			{}, // title is minOccurs 1, maxOccurs unbounded
		}
	}

	return descMeta.Title
}

// Appenders

func (typeDesc *TypeDesc) AppendDenomination(denomination Denomination) {
	denominations := typeDesc.Denomination
	if denominations == nil {
		denominations = []Denomination{}
	}

	typeDesc.Denomination = append(
		denominations,
		denomination)
}

func (typeDesc *TypeDesc) AppendMaterial(material Material) {
	materials := typeDesc.Material
	if materials == nil {
		materials = []Material{}
	}
	typeDesc.Material = append(
		materials,
		material)
}

func (fileGrp *FileGrp) AppendFile(file File) {
	files := fileGrp.File
	if files == nil {
		files = []File{}
	}
	fileGrp.File = append(
		files,
		file)
}

func (rightsStmt *RightsStmt) AppendLicense(license License) {
	licenses := rightsStmt.License
	if licenses == nil {
		licenses = []License{}
	}
	rightsStmt.License = append(
		licenses,
		license)
}

// Generators

func NewNUDS(recordType string) NUDS {
	return NUDS{
		XMLNS:          "http://nomisma.org/nuds",
		METS_NS:        "http://www.loc.gov/METS/",
		TEI_NS:         "http://www.tei-c.org/ns/1.0",
		XS_NS:          "http://www.w3.org/2001/XMLSchema",
		XLINK_NS:       "http://www.w3.org/1999/xlink",
		XSI_NS:         "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://nomisma.org/nuds http://nomisma.org/nuds.xsd",

		RecordType: recordType,

		Control: Control{
			PublicationStatus: PublicationStatus{
				Value: "inProcess",
			},
			MaintenanceStatus: MaintenanceStatus{
				Value: "derived",
			},
			MaintenanceHistory: MaintenanceHistory{
				MaintenanceEvent: []MaintenanceEvent{
					{
						EventType: EventType{Value: "derived"},
						EventDateTime: EventDateTime{
							StandardDateTime: time.Now().String(),
							Value:            time.Now().Format("01-02-2006 15:04:05"),
						},
						AgentType: AgentType{Value: "machine"},
						Agent:     Agent{"csv-nuds"},
					},
				},
			},
		},
	}
}
