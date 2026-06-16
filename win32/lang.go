package win32

type LANGID WORD

type LANGID_PRIMARY WORD
type LANGID_SUB WORD

type SORTID WORD

type LCID DWORD

const (
	LANG_CHINESE_TRADITIONAL LANGID = 0x7c04 // Use with the ConvertDefaultLocale function
	LANG_BOSNIAN_NEUTRAL     LANGID = 0x781a // Use with the ConvertDefaultLocale function
	LANG_SERBIAN_NEUTRAL     LANGID = 0x7c1a // Use with the ConvertDefaultLocale function

	LANG_NEUTRAL   LANGID_PRIMARY = 0x00
	LANG_INVARIANT LANGID_PRIMARY = 0x7f

	LANG_AFRIKAANS          LANGID_PRIMARY = 0x36
	LANG_ALBANIAN           LANGID_PRIMARY = 0x1c
	LANG_ALSATIAN           LANGID_PRIMARY = 0x84
	LANG_AMHARIC            LANGID_PRIMARY = 0x5e
	LANG_ARABIC             LANGID_PRIMARY = 0x01
	LANG_ARMENIAN           LANGID_PRIMARY = 0x2b
	LANG_ASSAMESE           LANGID_PRIMARY = 0x4d
	LANG_AZERI              LANGID_PRIMARY = 0x2c // for Azerbaijani, LANG_AZERBAIJANI is preferred
	LANG_AZERBAIJANI        LANGID_PRIMARY = 0x2c
	LANG_BANGLA             LANGID_PRIMARY = 0x45
	LANG_BASHKIR            LANGID_PRIMARY = 0x6d
	LANG_BASQUE             LANGID_PRIMARY = 0x2d
	LANG_BELARUSIAN         LANGID_PRIMARY = 0x23
	LANG_BENGALI            LANGID_PRIMARY = 0x45 // Some prefer to use LANG_BANGLA
	LANG_BRETON             LANGID_PRIMARY = 0x7e
	LANG_BOSNIAN            LANGID_PRIMARY = 0x1a // Use with SUBLANG_BOSNIAN_* Sublanguage IDs
	LANG_BULGARIAN          LANGID_PRIMARY = 0x02
	LANG_CATALAN            LANGID_PRIMARY = 0x03
	LANG_CENTRAL_KURDISH    LANGID_PRIMARY = 0x92
	LANG_CHEROKEE           LANGID_PRIMARY = 0x5c
	LANG_CHINESE            LANGID_PRIMARY = 0x04 // Use with SUBLANG_CHINESE_* Sublanguage IDs
	LANG_CHINESE_SIMPLIFIED LANGID_PRIMARY = 0x04 // Use with the ConvertDefaultLocale function
	LANG_CORSICAN           LANGID_PRIMARY = 0x83
	LANG_CROATIAN           LANGID_PRIMARY = 0x1a
	LANG_CZECH              LANGID_PRIMARY = 0x05
	LANG_DANISH             LANGID_PRIMARY = 0x06
	LANG_DARI               LANGID_PRIMARY = 0x8c
	LANG_DIVEHI             LANGID_PRIMARY = 0x65
	LANG_DUTCH              LANGID_PRIMARY = 0x13
	LANG_ENGLISH            LANGID_PRIMARY = 0x09
	LANG_ESTONIAN           LANGID_PRIMARY = 0x25
	LANG_FAEROESE           LANGID_PRIMARY = 0x38
	LANG_FARSI              LANGID_PRIMARY = 0x29 // Deprecated: use LANG_PERSIAN instead
	LANG_FILIPINO           LANGID_PRIMARY = 0x64
	LANG_FINNISH            LANGID_PRIMARY = 0x0b
	LANG_FRENCH             LANGID_PRIMARY = 0x0c
	LANG_FRISIAN            LANGID_PRIMARY = 0x62
	LANG_FULAH              LANGID_PRIMARY = 0x67
	LANG_GALICIAN           LANGID_PRIMARY = 0x56
	LANG_GEORGIAN           LANGID_PRIMARY = 0x37
	LANG_GERMAN             LANGID_PRIMARY = 0x07
	LANG_GREEK              LANGID_PRIMARY = 0x08
	LANG_GREENLANDIC        LANGID_PRIMARY = 0x6f
	LANG_GUJARATI           LANGID_PRIMARY = 0x47
	LANG_HAUSA              LANGID_PRIMARY = 0x68
	LANG_HAWAIIAN           LANGID_PRIMARY = 0x75
	LANG_HEBREW             LANGID_PRIMARY = 0x0d
	LANG_HINDI              LANGID_PRIMARY = 0x39
	LANG_HUNGARIAN          LANGID_PRIMARY = 0x0e
	LANG_ICELANDIC          LANGID_PRIMARY = 0x0f
	LANG_IGBO               LANGID_PRIMARY = 0x70
	LANG_INDONESIAN         LANGID_PRIMARY = 0x21
	LANG_INUKTITUT          LANGID_PRIMARY = 0x5d
	LANG_IRISH              LANGID_PRIMARY = 0x3c // Use with the SUBLANG_IRISH_IRELAND Sublanguage ID
	LANG_ITALIAN            LANGID_PRIMARY = 0x10
	LANG_JAPANESE           LANGID_PRIMARY = 0x11
	LANG_KANNADA            LANGID_PRIMARY = 0x4b
	LANG_KASHMIRI           LANGID_PRIMARY = 0x60
	LANG_KAZAK              LANGID_PRIMARY = 0x3f
	LANG_KHMER              LANGID_PRIMARY = 0x53
	LANG_KICHE              LANGID_PRIMARY = 0x86
	LANG_KINYARWANDA        LANGID_PRIMARY = 0x87
	LANG_KONKANI            LANGID_PRIMARY = 0x57
	LANG_KOREAN             LANGID_PRIMARY = 0x12
	LANG_KYRGYZ             LANGID_PRIMARY = 0x40
	LANG_LAO                LANGID_PRIMARY = 0x54
	LANG_LATVIAN            LANGID_PRIMARY = 0x26
	LANG_LITHUANIAN         LANGID_PRIMARY = 0x27
	LANG_LOWER_SORBIAN      LANGID_PRIMARY = 0x2e
	LANG_LUXEMBOURGISH      LANGID_PRIMARY = 0x6e
	LANG_MACEDONIAN         LANGID_PRIMARY = 0x2f // the Former Yugoslav Republic of Macedonia
	LANG_MALAY              LANGID_PRIMARY = 0x3e
	LANG_MALAYALAM          LANGID_PRIMARY = 0x4c
	LANG_MALTESE            LANGID_PRIMARY = 0x3a
	LANG_MANIPURI           LANGID_PRIMARY = 0x58
	LANG_MAORI              LANGID_PRIMARY = 0x81
	LANG_MAPUDUNGUN         LANGID_PRIMARY = 0x7a
	LANG_MARATHI            LANGID_PRIMARY = 0x4e
	LANG_MOHAWK             LANGID_PRIMARY = 0x7c
	LANG_MONGOLIAN          LANGID_PRIMARY = 0x50
	LANG_NEPALI             LANGID_PRIMARY = 0x61
	LANG_NORWEGIAN          LANGID_PRIMARY = 0x14
	LANG_OCCITAN            LANGID_PRIMARY = 0x82
	LANG_ODIA               LANGID_PRIMARY = 0x48
	LANG_ORIYA              LANGID_PRIMARY = 0x48 // Deprecated: use LANG_ODIA, instead.
	LANG_PASHTO             LANGID_PRIMARY = 0x63
	LANG_PERSIAN            LANGID_PRIMARY = 0x29
	LANG_POLISH             LANGID_PRIMARY = 0x15
	LANG_PORTUGUESE         LANGID_PRIMARY = 0x16
	LANG_PULAR              LANGID_PRIMARY = 0x67 // Deprecated: use LANG_FULAH instead
	LANG_PUNJABI            LANGID_PRIMARY = 0x46
	LANG_QUECHUA            LANGID_PRIMARY = 0x6b
	LANG_ROMANIAN           LANGID_PRIMARY = 0x18
	LANG_ROMANSH            LANGID_PRIMARY = 0x17
	LANG_RUSSIAN            LANGID_PRIMARY = 0x19
	LANG_SAKHA              LANGID_PRIMARY = 0x85
	LANG_SAMI               LANGID_PRIMARY = 0x3b
	LANG_SANSKRIT           LANGID_PRIMARY = 0x4f
	LANG_SCOTTISH_GAELIC    LANGID_PRIMARY = 0x91
	LANG_SERBIAN            LANGID_PRIMARY = 0x1a // Use with the SUBLANG_SERBIAN_* Sublanguage IDs
	LANG_SINDHI             LANGID_PRIMARY = 0x59
	LANG_SINHALESE          LANGID_PRIMARY = 0x5b
	LANG_SLOVAK             LANGID_PRIMARY = 0x1b
	LANG_SLOVENIAN          LANGID_PRIMARY = 0x24
	LANG_SOTHO              LANGID_PRIMARY = 0x6c
	LANG_SPANISH            LANGID_PRIMARY = 0x0a
	LANG_SWAHILI            LANGID_PRIMARY = 0x41
	LANG_SWEDISH            LANGID_PRIMARY = 0x1d
	LANG_SYRIAC             LANGID_PRIMARY = 0x5a
	LANG_TAJIK              LANGID_PRIMARY = 0x28
	LANG_TAMAZIGHT          LANGID_PRIMARY = 0x5f
	LANG_TAMIL              LANGID_PRIMARY = 0x49
	LANG_TATAR              LANGID_PRIMARY = 0x44
	LANG_TELUGU             LANGID_PRIMARY = 0x4a
	LANG_THAI               LANGID_PRIMARY = 0x1e
	LANG_TIBETAN            LANGID_PRIMARY = 0x51
	LANG_TIGRIGNA           LANGID_PRIMARY = 0x73
	LANG_TIGRINYA           LANGID_PRIMARY = 0x73 // Preferred spelling in locale
	LANG_TSWANA             LANGID_PRIMARY = 0x32
	LANG_TURKISH            LANGID_PRIMARY = 0x1f
	LANG_TURKMEN            LANGID_PRIMARY = 0x42
	LANG_UIGHUR             LANGID_PRIMARY = 0x80
	LANG_UKRAINIAN          LANGID_PRIMARY = 0x22
	LANG_UPPER_SORBIAN      LANGID_PRIMARY = 0x2e
	LANG_URDU               LANGID_PRIMARY = 0x20
	LANG_UZBEK              LANGID_PRIMARY = 0x43
	LANG_VALENCIAN          LANGID_PRIMARY = 0x03
	LANG_VIETNAMESE         LANGID_PRIMARY = 0x2a
	LANG_WELSH              LANGID_PRIMARY = 0x52
	LANG_WOLOF              LANGID_PRIMARY = 0x88
	LANG_XHOSA              LANGID_PRIMARY = 0x34
	LANG_YAKUT              LANGID_PRIMARY = 0x85 // Deprecated: use LANG_SAKHA,instead
	LANG_YI                 LANGID_PRIMARY = 0x78
	LANG_YORUBA             LANGID_PRIMARY = 0x6a
	LANG_ZULU               LANGID_PRIMARY = 0x35

	SUBLANG_NEUTRAL            LANGID_SUB = 0x00 // language neutral
	SUBLANG_DEFAULT            LANGID_SUB = 0x01 // user default
	SUBLANG_SYS_DEFAULT        LANGID_SUB = 0x02 // system default
	SUBLANG_CUSTOM_DEFAULT     LANGID_SUB = 0x03 // default custom language/locale
	SUBLANG_CUSTOM_UNSPECIFIED LANGID_SUB = 0x04 // custom language/locale
	SUBLANG_UI_CUSTOM_DEFAULT  LANGID_SUB = 0x05 // Default custom MUI language/locale

	SUBLANG_AFRIKAANS_SOUTH_AFRICA              LANGID_SUB = 0x01 // Afrikaans (South Africa) 0x0436 af-ZA
	SUBLANG_ALBANIAN_ALBANIA                    LANGID_SUB = 0x01 // Albanian (Albania) 0x041c sq-AL
	SUBLANG_ALSATIAN_FRANCE                     LANGID_SUB = 0x01 // Alsatian (France) 0x0484
	SUBLANG_AMHARIC_ETHIOPIA                    LANGID_SUB = 0x01 // Amharic (Ethiopia) 0x045e
	SUBLANG_ARABIC_SAUDI_ARABIA                 LANGID_SUB = 0x01 // Arabic (Saudi Arabia)
	SUBLANG_ARABIC_IRAQ                         LANGID_SUB = 0x02 // Arabic (Iraq)
	SUBLANG_ARABIC_EGYPT                        LANGID_SUB = 0x03 // Arabic (Egypt)
	SUBLANG_ARABIC_LIBYA                        LANGID_SUB = 0x04 // Arabic (Libya)
	SUBLANG_ARABIC_ALGERIA                      LANGID_SUB = 0x05 // Arabic (Algeria)
	SUBLANG_ARABIC_MOROCCO                      LANGID_SUB = 0x06 // Arabic (Morocco)
	SUBLANG_ARABIC_TUNISIA                      LANGID_SUB = 0x07 // Arabic (Tunisia)
	SUBLANG_ARABIC_OMAN                         LANGID_SUB = 0x08 // Arabic (Oman)
	SUBLANG_ARABIC_YEMEN                        LANGID_SUB = 0x09 // Arabic (Yemen)
	SUBLANG_ARABIC_SYRIA                        LANGID_SUB = 0x0a // Arabic (Syria)
	SUBLANG_ARABIC_JORDAN                       LANGID_SUB = 0x0b // Arabic (Jordan)
	SUBLANG_ARABIC_LEBANON                      LANGID_SUB = 0x0c // Arabic (Lebanon)
	SUBLANG_ARABIC_KUWAIT                       LANGID_SUB = 0x0d // Arabic (Kuwait)
	SUBLANG_ARABIC_UAE                          LANGID_SUB = 0x0e // Arabic (U.A.E)
	SUBLANG_ARABIC_BAHRAIN                      LANGID_SUB = 0x0f // Arabic (Bahrain)
	SUBLANG_ARABIC_QATAR                        LANGID_SUB = 0x10 // Arabic (Qatar)
	SUBLANG_ARMENIAN_ARMENIA                    LANGID_SUB = 0x01 // Armenian (Armenia) 0x042b hy-AM
	SUBLANG_ASSAMESE_INDIA                      LANGID_SUB = 0x01 // Assamese (India) 0x044d
	SUBLANG_AZERI_LATIN                         LANGID_SUB = 0x01 // Azeri (Latin) - for Azerbaijani, SUBLANG_AZERBAIJANI_AZERBAIJAN_LATIN preferred
	SUBLANG_AZERI_CYRILLIC                      LANGID_SUB = 0x02 // Azeri (Cyrillic) - for Azerbaijani, SUBLANG_AZERBAIJANI_AZERBAIJAN_CYRILLIC preferred
	SUBLANG_AZERBAIJANI_AZERBAIJAN_LATIN        LANGID_SUB = 0x01 // Azerbaijani (Azerbaijan, Latin)
	SUBLANG_AZERBAIJANI_AZERBAIJAN_CYRILLIC     LANGID_SUB = 0x02 // Azerbaijani (Azerbaijan, Cyrillic)
	SUBLANG_BANGLA_INDIA                        LANGID_SUB = 0x01 // Bangla (India)
	SUBLANG_BANGLA_BANGLADESH                   LANGID_SUB = 0x02 // Bangla (Bangladesh)
	SUBLANG_BASHKIR_RUSSIA                      LANGID_SUB = 0x01 // Bashkir (Russia) 0x046d ba-RU
	SUBLANG_BASQUE_BASQUE                       LANGID_SUB = 0x01 // Basque (Basque) 0x042d eu-ES
	SUBLANG_BELARUSIAN_BELARUS                  LANGID_SUB = 0x01 // Belarusian (Belarus) 0x0423 be-BY
	SUBLANG_BENGALI_INDIA                       LANGID_SUB = 0x01 // Bengali (India) - Note some prefer SUBLANG_BANGLA_INDIA
	SUBLANG_BENGALI_BANGLADESH                  LANGID_SUB = 0x02 // Bengali (Bangladesh) - Note some prefer SUBLANG_BANGLA_BANGLADESH
	SUBLANG_BOSNIAN_BOSNIA_HERZEGOVINA_LATIN    LANGID_SUB = 0x05 // Bosnian (Bosnia and Herzegovina - Latin) 0x141a bs-BA-Latn
	SUBLANG_BOSNIAN_BOSNIA_HERZEGOVINA_CYRILLIC LANGID_SUB = 0x08 // Bosnian (Bosnia and Herzegovina - Cyrillic) 0x201a bs-BA-Cyrl
	SUBLANG_BRETON_FRANCE                       LANGID_SUB = 0x01 // Breton (France) 0x047e
	SUBLANG_BULGARIAN_BULGARIA                  LANGID_SUB = 0x01 // Bulgarian (Bulgaria) 0x0402
	SUBLANG_CATALAN_CATALAN                     LANGID_SUB = 0x01 // Catalan (Catalan) 0x0403
	SUBLANG_CENTRAL_KURDISH_IRAQ                LANGID_SUB = 0x01 // Central Kurdish (Iraq) 0x0492 ku-Arab-IQ
	SUBLANG_CHEROKEE_CHEROKEE                   LANGID_SUB = 0x01 // Cherokee (Cherokee) 0x045c chr-Cher-US
	SUBLANG_CHINESE_TRADITIONAL                 LANGID_SUB = 0x01 // Chinese (Taiwan) 0x0404 zh-TW
	SUBLANG_CHINESE_SIMPLIFIED                  LANGID_SUB = 0x02 // Chinese (PR China) 0x0804 zh-CN
	SUBLANG_CHINESE_HONGKONG                    LANGID_SUB = 0x03 // Chinese (Hong Kong S.A.R., P.R.C.) 0x0c04 zh-HK
	SUBLANG_CHINESE_SINGAPORE                   LANGID_SUB = 0x04 // Chinese (Singapore) 0x1004 zh-SG
	SUBLANG_CHINESE_MACAU                       LANGID_SUB = 0x05 // Chinese (Macau S.A.R.) 0x1404 zh-MO
	SUBLANG_CORSICAN_FRANCE                     LANGID_SUB = 0x01 // Corsican (France) 0x0483
	SUBLANG_CZECH_CZECH_REPUBLIC                LANGID_SUB = 0x01 // Czech (Czech Republic) 0x0405
	SUBLANG_CROATIAN_CROATIA                    LANGID_SUB = 0x01 // Croatian (Croatia)
	SUBLANG_CROATIAN_BOSNIA_HERZEGOVINA_LATIN   LANGID_SUB = 0x04 // Croatian (Bosnia and Herzegovina - Latin) 0x101a hr-BA
	SUBLANG_DANISH_DENMARK                      LANGID_SUB = 0x01 // Danish (Denmark) 0x0406
	SUBLANG_DARI_AFGHANISTAN                    LANGID_SUB = 0x01 // Dari (Afghanistan)
	SUBLANG_DIVEHI_MALDIVES                     LANGID_SUB = 0x01 // Divehi (Maldives) 0x0465 div-MV
	SUBLANG_DUTCH                               LANGID_SUB = 0x01 // Dutch
	SUBLANG_DUTCH_BELGIAN                       LANGID_SUB = 0x02 // Dutch (Belgian)
	SUBLANG_ENGLISH_US                          LANGID_SUB = 0x01 // English (USA)
	SUBLANG_ENGLISH_UK                          LANGID_SUB = 0x02 // English (UK)
	SUBLANG_ENGLISH_AUS                         LANGID_SUB = 0x03 // English (Australian)
	SUBLANG_ENGLISH_CAN                         LANGID_SUB = 0x04 // English (Canadian)
	SUBLANG_ENGLISH_NZ                          LANGID_SUB = 0x05 // English (New Zealand)
	SUBLANG_ENGLISH_EIRE                        LANGID_SUB = 0x06 // English (Irish)
	SUBLANG_ENGLISH_SOUTH_AFRICA                LANGID_SUB = 0x07 // English (South Africa)
	SUBLANG_ENGLISH_JAMAICA                     LANGID_SUB = 0x08 // English (Jamaica)
	SUBLANG_ENGLISH_CARIBBEAN                   LANGID_SUB = 0x09 // English (Caribbean)
	SUBLANG_ENGLISH_BELIZE                      LANGID_SUB = 0x0a // English (Belize)
	SUBLANG_ENGLISH_TRINIDAD                    LANGID_SUB = 0x0b // English (Trinidad)
	SUBLANG_ENGLISH_ZIMBABWE                    LANGID_SUB = 0x0c // English (Zimbabwe)
	SUBLANG_ENGLISH_PHILIPPINES                 LANGID_SUB = 0x0d // English (Philippines)
	SUBLANG_ENGLISH_INDIA                       LANGID_SUB = 0x10 // English (India)
	SUBLANG_ENGLISH_MALAYSIA                    LANGID_SUB = 0x11 // English (Malaysia)
	SUBLANG_ENGLISH_SINGAPORE                   LANGID_SUB = 0x12 // English (Singapore)
	SUBLANG_ESTONIAN_ESTONIA                    LANGID_SUB = 0x01 // Estonian (Estonia) 0x0425 et-EE
	SUBLANG_FAEROESE_FAROE_ISLANDS              LANGID_SUB = 0x01 // Faroese (Faroe Islands) 0x0438 fo-FO
	SUBLANG_FILIPINO_PHILIPPINES                LANGID_SUB = 0x01 // Filipino (Philippines) 0x0464 fil-PH
	SUBLANG_FINNISH_FINLAND                     LANGID_SUB = 0x01 // Finnish (Finland) 0x040b
	SUBLANG_FRENCH                              LANGID_SUB = 0x01 // French
	SUBLANG_FRENCH_BELGIAN                      LANGID_SUB = 0x02 // French (Belgian)
	SUBLANG_FRENCH_CANADIAN                     LANGID_SUB = 0x03 // French (Canadian)
	SUBLANG_FRENCH_SWISS                        LANGID_SUB = 0x04 // French (Swiss)
	SUBLANG_FRENCH_LUXEMBOURG                   LANGID_SUB = 0x05 // French (Luxembourg)
	SUBLANG_FRENCH_MONACO                       LANGID_SUB = 0x06 // French (Monaco)
	SUBLANG_FRISIAN_NETHERLANDS                 LANGID_SUB = 0x01 // Frisian (Netherlands) 0x0462 fy-NL
	SUBLANG_FULAH_SENEGAL                       LANGID_SUB = 0x02 // Fulah (Senegal) 0x0867 ff-Latn-SN
	SUBLANG_GALICIAN_GALICIAN                   LANGID_SUB = 0x01 // Galician (Galician) 0x0456 gl-ES
	SUBLANG_GEORGIAN_GEORGIA                    LANGID_SUB = 0x01 // Georgian (Georgia) 0x0437 ka-GE
	SUBLANG_GERMAN                              LANGID_SUB = 0x01 // German
	SUBLANG_GERMAN_SWISS                        LANGID_SUB = 0x02 // German (Swiss)
	SUBLANG_GERMAN_AUSTRIAN                     LANGID_SUB = 0x03 // German (Austrian)
	SUBLANG_GERMAN_LUXEMBOURG                   LANGID_SUB = 0x04 // German (Luxembourg)
	SUBLANG_GERMAN_LIECHTENSTEIN                LANGID_SUB = 0x05 // German (Liechtenstein)
	SUBLANG_GREEK_GREECE                        LANGID_SUB = 0x01 // Greek (Greece)
	SUBLANG_GREENLANDIC_GREENLAND               LANGID_SUB = 0x01 // Greenlandic (Greenland) 0x046f kl-GL
	SUBLANG_GUJARATI_INDIA                      LANGID_SUB = 0x01 // Gujarati (India (Gujarati Script)) 0x0447 gu-IN
	SUBLANG_HAUSA_NIGERIA_LATIN                 LANGID_SUB = 0x01 // Hausa (Latin, Nigeria) 0x0468 ha-NG-Latn
	SUBLANG_HAWAIIAN_US                         LANGID_SUB = 0x01 // Hawiian (US) 0x0475 haw-US
	SUBLANG_HEBREW_ISRAEL                       LANGID_SUB = 0x01 // Hebrew (Israel) 0x040d
	SUBLANG_HINDI_INDIA                         LANGID_SUB = 0x01 // Hindi (India) 0x0439 hi-IN
	SUBLANG_HUNGARIAN_HUNGARY                   LANGID_SUB = 0x01 // Hungarian (Hungary) 0x040e
	SUBLANG_ICELANDIC_ICELAND                   LANGID_SUB = 0x01 // Icelandic (Iceland) 0x040f
	SUBLANG_IGBO_NIGERIA                        LANGID_SUB = 0x01 // Igbo (Nigeria) 0x0470 ig-NG
	SUBLANG_INDONESIAN_INDONESIA                LANGID_SUB = 0x01 // Indonesian (Indonesia) 0x0421 id-ID
	SUBLANG_INUKTITUT_CANADA                    LANGID_SUB = 0x01 // Inuktitut (Syllabics) (Canada) 0x045d iu-CA-Cans
	SUBLANG_INUKTITUT_CANADA_LATIN              LANGID_SUB = 0x02 // Inuktitut (Canada - Latin)
	SUBLANG_IRISH_IRELAND                       LANGID_SUB = 0x02 // Irish (Ireland)
	SUBLANG_ITALIAN                             LANGID_SUB = 0x01 // Italian
	SUBLANG_ITALIAN_SWISS                       LANGID_SUB = 0x02 // Italian (Swiss)
	SUBLANG_JAPANESE_JAPAN                      LANGID_SUB = 0x01 // Japanese (Japan) 0x0411
	SUBLANG_KANNADA_INDIA                       LANGID_SUB = 0x01 // Kannada (India (Kannada Script)) 0x044b kn-IN
	SUBLANG_KASHMIRI_SASIA                      LANGID_SUB = 0x02 // Kashmiri (South Asia)
	SUBLANG_KASHMIRI_INDIA                      LANGID_SUB = 0x02 // For app compatibility only
	SUBLANG_KAZAK_KAZAKHSTAN                    LANGID_SUB = 0x01 // Kazakh (Kazakhstan) 0x043f kk-KZ
	SUBLANG_KHMER_CAMBODIA                      LANGID_SUB = 0x01 // Khmer (Cambodia) 0x0453 kh-KH
	SUBLANG_KICHE_GUATEMALA                     LANGID_SUB = 0x01 // K'iche (Guatemala)
	SUBLANG_KINYARWANDA_RWANDA                  LANGID_SUB = 0x01 // Kinyarwanda (Rwanda) 0x0487 rw-RW
	SUBLANG_KONKANI_INDIA                       LANGID_SUB = 0x01 // Konkani (India) 0x0457 kok-IN
	SUBLANG_KOREAN                              LANGID_SUB = 0x01 // Korean (Extended Wansung)
	SUBLANG_KYRGYZ_KYRGYZSTAN                   LANGID_SUB = 0x01 // Kyrgyz (Kyrgyzstan) 0x0440 ky-KG
	SUBLANG_LAO_LAO                             LANGID_SUB = 0x01 // Lao (Lao PDR) 0x0454 lo-LA
	SUBLANG_LATVIAN_LATVIA                      LANGID_SUB = 0x01 // Latvian (Latvia) 0x0426 lv-LV
	SUBLANG_LITHUANIAN                          LANGID_SUB = 0x01 // Lithuanian
	SUBLANG_LOWER_SORBIAN_GERMANY               LANGID_SUB = 0x02 // Lower Sorbian (Germany) 0x082e wee-DE
	SUBLANG_LUXEMBOURGISH_LUXEMBOURG            LANGID_SUB = 0x01 // Luxembourgish (Luxembourg) 0x046e lb-LU
	SUBLANG_MACEDONIAN_MACEDONIA                LANGID_SUB = 0x01 // Macedonian (Macedonia (FYROM)) 0x042f mk-MK
	SUBLANG_MALAY_MALAYSIA                      LANGID_SUB = 0x01 // Malay (Malaysia)
	SUBLANG_MALAY_BRUNEI_DARUSSALAM             LANGID_SUB = 0x02 // Malay (Brunei Darussalam)
	SUBLANG_MALAYALAM_INDIA                     LANGID_SUB = 0x01 // Malayalam (India (Malayalam Script) ) 0x044c ml-IN
	SUBLANG_MALTESE_MALTA                       LANGID_SUB = 0x01 // Maltese (Malta) 0x043a mt-MT
	SUBLANG_MAORI_NEW_ZEALAND                   LANGID_SUB = 0x01 // Maori (New Zealand) 0x0481 mi-NZ
	SUBLANG_MAPUDUNGUN_CHILE                    LANGID_SUB = 0x01 // Mapudungun (Chile) 0x047a arn-CL
	SUBLANG_MARATHI_INDIA                       LANGID_SUB = 0x01 // Marathi (India) 0x044e mr-IN
	SUBLANG_MOHAWK_MOHAWK                       LANGID_SUB = 0x01 // Mohawk (Mohawk) 0x047c moh-CA
	SUBLANG_MONGOLIAN_CYRILLIC_MONGOLIA         LANGID_SUB = 0x01 // Mongolian (Cyrillic, Mongolia)
	SUBLANG_MONGOLIAN_PRC                       LANGID_SUB = 0x02 // Mongolian (PRC)
	SUBLANG_NEPALI_INDIA                        LANGID_SUB = 0x02 // Nepali (India)
	SUBLANG_NEPALI_NEPAL                        LANGID_SUB = 0x01 // Nepali (Nepal) 0x0461 ne-NP
	SUBLANG_NORWEGIAN_BOKMAL                    LANGID_SUB = 0x01 // Norwegian (Bokmal)
	SUBLANG_NORWEGIAN_NYNORSK                   LANGID_SUB = 0x02 // Norwegian (Nynorsk)
	SUBLANG_OCCITAN_FRANCE                      LANGID_SUB = 0x01 // Occitan (France) 0x0482 oc-FR
	SUBLANG_ODIA_INDIA                          LANGID_SUB = 0x01 // Odia (India (Odia Script)) 0x0448 or-IN
	SUBLANG_ORIYA_INDIA                         LANGID_SUB = 0x01 // Deprecated: use SUBLANG_ODIA_INDIA instead
	SUBLANG_PASHTO_AFGHANISTAN                  LANGID_SUB = 0x01 // Pashto (Afghanistan)
	SUBLANG_PERSIAN_IRAN                        LANGID_SUB = 0x01 // Persian (Iran) 0x0429 fa-IR
	SUBLANG_POLISH_POLAND                       LANGID_SUB = 0x01 // Polish (Poland) 0x0415
	SUBLANG_PORTUGUESE                          LANGID_SUB = 0x02 // Portuguese
	SUBLANG_PORTUGUESE_BRAZILIAN                LANGID_SUB = 0x01 // Portuguese (Brazil)
	SUBLANG_PULAR_SENEGAL                       LANGID_SUB = 0x02 // Deprecated: Use SUBLANG_FULAH_SENEGAL instead
	SUBLANG_PUNJABI_INDIA                       LANGID_SUB = 0x01 // Punjabi (India (Gurmukhi Script)) 0x0446 pa-IN
	SUBLANG_PUNJABI_PAKISTAN                    LANGID_SUB = 0x02 // Punjabi (Pakistan (Arabic Script)) 0x0846 pa-Arab-PK
	SUBLANG_QUECHUA_BOLIVIA                     LANGID_SUB = 0x01 // Quechua (Bolivia)
	SUBLANG_QUECHUA_ECUADOR                     LANGID_SUB = 0x02 // Quechua (Ecuador)
	SUBLANG_QUECHUA_PERU                        LANGID_SUB = 0x03 // Quechua (Peru)
	SUBLANG_ROMANIAN_ROMANIA                    LANGID_SUB = 0x01 // Romanian (Romania) 0x0418
	SUBLANG_ROMANSH_SWITZERLAND                 LANGID_SUB = 0x01 // Romansh (Switzerland) 0x0417 rm-CH
	SUBLANG_RUSSIAN_RUSSIA                      LANGID_SUB = 0x01 // Russian (Russia) 0x0419
	SUBLANG_SAKHA_RUSSIA                        LANGID_SUB = 0x01 // Sakha (Russia) 0x0485 sah-RU
	SUBLANG_SAMI_NORTHERN_NORWAY                LANGID_SUB = 0x01 // Northern Sami (Norway)
	SUBLANG_SAMI_NORTHERN_SWEDEN                LANGID_SUB = 0x02 // Northern Sami (Sweden)
	SUBLANG_SAMI_NORTHERN_FINLAND               LANGID_SUB = 0x03 // Northern Sami (Finland)
	SUBLANG_SAMI_LULE_NORWAY                    LANGID_SUB = 0x04 // Lule Sami (Norway)
	SUBLANG_SAMI_LULE_SWEDEN                    LANGID_SUB = 0x05 // Lule Sami (Sweden)
	SUBLANG_SAMI_SOUTHERN_NORWAY                LANGID_SUB = 0x06 // Southern Sami (Norway)
	SUBLANG_SAMI_SOUTHERN_SWEDEN                LANGID_SUB = 0x07 // Southern Sami (Sweden)
	SUBLANG_SAMI_SKOLT_FINLAND                  LANGID_SUB = 0x08 // Skolt Sami (Finland)
	SUBLANG_SAMI_INARI_FINLAND                  LANGID_SUB = 0x09 // Inari Sami (Finland)
	SUBLANG_SANSKRIT_INDIA                      LANGID_SUB = 0x01 // Sanskrit (India) 0x044f sa-IN
	SUBLANG_SCOTTISH_GAELIC                     LANGID_SUB = 0x01 // Scottish Gaelic (United Kingdom) 0x0491 gd-GB
	SUBLANG_SERBIAN_BOSNIA_HERZEGOVINA_LATIN    LANGID_SUB = 0x06 // Serbian (Bosnia and Herzegovina - Latin)
	SUBLANG_SERBIAN_BOSNIA_HERZEGOVINA_CYRILLIC LANGID_SUB = 0x07 // Serbian (Bosnia and Herzegovina - Cyrillic)
	SUBLANG_SERBIAN_MONTENEGRO_LATIN            LANGID_SUB = 0x0b // Serbian (Montenegro - Latn)
	SUBLANG_SERBIAN_MONTENEGRO_CYRILLIC         LANGID_SUB = 0x0c // Serbian (Montenegro - Cyrillic)
	SUBLANG_SERBIAN_SERBIA_LATIN                LANGID_SUB = 0x09 // Serbian (Serbia - Latin)
	SUBLANG_SERBIAN_SERBIA_CYRILLIC             LANGID_SUB = 0x0a // Serbian (Serbia - Cyrillic)
	SUBLANG_SERBIAN_CROATIA                     LANGID_SUB = 0x01 // Croatian (Croatia) 0x041a hr-HR
	SUBLANG_SERBIAN_LATIN                       LANGID_SUB = 0x02 // Serbian (Latin)
	SUBLANG_SERBIAN_CYRILLIC                    LANGID_SUB = 0x03 // Serbian (Cyrillic)
	SUBLANG_SINDHI_INDIA                        LANGID_SUB = 0x01 // Sindhi (India) reserved 0x0459
	SUBLANG_SINDHI_PAKISTAN                     LANGID_SUB = 0x02 // Sindhi (Pakistan) 0x0859 sd-Arab-PK
	SUBLANG_SINDHI_AFGHANISTAN                  LANGID_SUB = 0x02 // For app compatibility only
	SUBLANG_SINHALESE_SRI_LANKA                 LANGID_SUB = 0x01 // Sinhalese (Sri Lanka)
	SUBLANG_SOTHO_NORTHERN_SOUTH_AFRICA         LANGID_SUB = 0x01 // Northern Sotho (South Africa)
	SUBLANG_SLOVAK_SLOVAKIA                     LANGID_SUB = 0x01 // Slovak (Slovakia) 0x041b sk-SK
	SUBLANG_SLOVENIAN_SLOVENIA                  LANGID_SUB = 0x01 // Slovenian (Slovenia) 0x0424 sl-SI
	SUBLANG_SPANISH                             LANGID_SUB = 0x01 // Spanish (Castilian)
	SUBLANG_SPANISH_MEXICAN                     LANGID_SUB = 0x02 // Spanish (Mexico)
	SUBLANG_SPANISH_MODERN                      LANGID_SUB = 0x03 // Spanish (Modern)
	SUBLANG_SPANISH_GUATEMALA                   LANGID_SUB = 0x04 // Spanish (Guatemala)
	SUBLANG_SPANISH_COSTA_RICA                  LANGID_SUB = 0x05 // Spanish (Costa Rica)
	SUBLANG_SPANISH_PANAMA                      LANGID_SUB = 0x06 // Spanish (Panama)
	SUBLANG_SPANISH_DOMINICAN_REPUBLIC          LANGID_SUB = 0x07 // Spanish (Dominican Republic)
	SUBLANG_SPANISH_VENEZUELA                   LANGID_SUB = 0x08 // Spanish (Venezuela)
	SUBLANG_SPANISH_COLOMBIA                    LANGID_SUB = 0x09 // Spanish (Colombia)
	SUBLANG_SPANISH_PERU                        LANGID_SUB = 0x0a // Spanish (Peru)
	SUBLANG_SPANISH_ARGENTINA                   LANGID_SUB = 0x0b // Spanish (Argentina)
	SUBLANG_SPANISH_ECUADOR                     LANGID_SUB = 0x0c // Spanish (Ecuador)
	SUBLANG_SPANISH_CHILE                       LANGID_SUB = 0x0d // Spanish (Chile)
	SUBLANG_SPANISH_URUGUAY                     LANGID_SUB = 0x0e // Spanish (Uruguay)
	SUBLANG_SPANISH_PARAGUAY                    LANGID_SUB = 0x0f // Spanish (Paraguay)
	SUBLANG_SPANISH_BOLIVIA                     LANGID_SUB = 0x10 // Spanish (Bolivia)
	SUBLANG_SPANISH_EL_SALVADOR                 LANGID_SUB = 0x11 // Spanish (El Salvador)
	SUBLANG_SPANISH_HONDURAS                    LANGID_SUB = 0x12 // Spanish (Honduras)
	SUBLANG_SPANISH_NICARAGUA                   LANGID_SUB = 0x13 // Spanish (Nicaragua)
	SUBLANG_SPANISH_PUERTO_RICO                 LANGID_SUB = 0x14 // Spanish (Puerto Rico)
	SUBLANG_SPANISH_US                          LANGID_SUB = 0x15 // Spanish (United States)
	SUBLANG_SWAHILI_KENYA                       LANGID_SUB = 0x01 // Swahili (Kenya) 0x0441 sw-KE
	SUBLANG_SWEDISH                             LANGID_SUB = 0x01 // Swedish
	SUBLANG_SWEDISH_FINLAND                     LANGID_SUB = 0x02 // Swedish (Finland)
	SUBLANG_SYRIAC_SYRIA                        LANGID_SUB = 0x01 // Syriac (Syria) 0x045a syr-SY
	SUBLANG_TAJIK_TAJIKISTAN                    LANGID_SUB = 0x01 // Tajik (Tajikistan) 0x0428 tg-TJ-Cyrl
	SUBLANG_TAMAZIGHT_ALGERIA_LATIN             LANGID_SUB = 0x02 // Tamazight (Latin, Algeria) 0x085f tzm-Latn-DZ
	SUBLANG_TAMAZIGHT_MOROCCO_TIFINAGH          LANGID_SUB = 0x04 // Tamazight (Tifinagh) 0x105f tzm-Tfng-MA
	SUBLANG_TAMIL_INDIA                         LANGID_SUB = 0x01 // Tamil (India)
	SUBLANG_TAMIL_SRI_LANKA                     LANGID_SUB = 0x02 // Tamil (Sri Lanka) 0x0849 ta-LK
	SUBLANG_TATAR_RUSSIA                        LANGID_SUB = 0x01 // Tatar (Russia) 0x0444 tt-RU
	SUBLANG_TELUGU_INDIA                        LANGID_SUB = 0x01 // Telugu (India (Telugu Script)) 0x044a te-IN
	SUBLANG_THAI_THAILAND                       LANGID_SUB = 0x01 // Thai (Thailand) 0x041e th-TH
	SUBLANG_TIBETAN_PRC                         LANGID_SUB = 0x01 // Tibetan (PRC)
	SUBLANG_TIGRIGNA_ERITREA                    LANGID_SUB = 0x02 // Tigrigna (Eritrea)
	SUBLANG_TIGRINYA_ERITREA                    LANGID_SUB = 0x02 // Tigrinya (Eritrea) 0x0873 ti-ER (preferred spelling)
	SUBLANG_TIGRINYA_ETHIOPIA                   LANGID_SUB = 0x01 // Tigrinya (Ethiopia) 0x0473 ti-ET
	SUBLANG_TSWANA_BOTSWANA                     LANGID_SUB = 0x02 // Setswana / Tswana (Botswana) 0x0832 tn-BW
	SUBLANG_TSWANA_SOUTH_AFRICA                 LANGID_SUB = 0x01 // Setswana / Tswana (South Africa) 0x0432 tn-ZA
	SUBLANG_TURKISH_TURKEY                      LANGID_SUB = 0x01 // Turkish (Turkey) 0x041f tr-TR
	SUBLANG_TURKMEN_TURKMENISTAN                LANGID_SUB = 0x01 // Turkmen (Turkmenistan) 0x0442 tk-TM
	SUBLANG_UIGHUR_PRC                          LANGID_SUB = 0x01 // Uighur (PRC) 0x0480 ug-CN
	SUBLANG_UKRAINIAN_UKRAINE                   LANGID_SUB = 0x01 // Ukrainian (Ukraine) 0x0422 uk-UA
	SUBLANG_UPPER_SORBIAN_GERMANY               LANGID_SUB = 0x01 // Upper Sorbian (Germany) 0x042e wen-DE
	SUBLANG_URDU_PAKISTAN                       LANGID_SUB = 0x01 // Urdu (Pakistan)
	SUBLANG_URDU_INDIA                          LANGID_SUB = 0x02 // Urdu (India)
	SUBLANG_UZBEK_LATIN                         LANGID_SUB = 0x01 // Uzbek (Latin)
	SUBLANG_UZBEK_CYRILLIC                      LANGID_SUB = 0x02 // Uzbek (Cyrillic)
	SUBLANG_VALENCIAN_VALENCIA                  LANGID_SUB = 0x02 // Valencian (Valencia) 0x0803 ca-ES-Valencia
	SUBLANG_VIETNAMESE_VIETNAM                  LANGID_SUB = 0x01 // Vietnamese (Vietnam) 0x042a vi-VN
	SUBLANG_WELSH_UNITED_KINGDOM                LANGID_SUB = 0x01 // Welsh (United Kingdom) 0x0452 cy-GB
	SUBLANG_WOLOF_SENEGAL                       LANGID_SUB = 0x01 // Wolof (Senegal)
	SUBLANG_XHOSA_SOUTH_AFRICA                  LANGID_SUB = 0x01 // isiXhosa / Xhosa (South Africa) 0x0434 xh-ZA
	SUBLANG_YAKUT_RUSSIA                        LANGID_SUB = 0x01 // Deprecated: use SUBLANG_SAKHA_RUSSIA instead
	SUBLANG_YI_PRC                              LANGID_SUB = 0x01 // Yi (PRC)) 0x0478
	SUBLANG_YORUBA_NIGERIA                      LANGID_SUB = 0x01 // Yoruba (Nigeria) 046a yo-NG
	SUBLANG_ZULU_SOUTH_AFRICA                   LANGID_SUB = 0x01 // isiZulu / Zulu (South Africa) 0x0435 zu-ZA

	SORT_DEFAULT                SORTID = 0x0 // sorting default
	SORT_INVARIANT_MATH         SORTID = 0x1 // Invariant (Mathematical Symbols)
	SORT_JAPANESE_XJIS          SORTID = 0x0 // Japanese XJIS order
	SORT_JAPANESE_UNICODE       SORTID = 0x1 // Japanese Unicode order (no longer supported)
	SORT_JAPANESE_RADICALSTROKE SORTID = 0x4 // Japanese radical/stroke order
	SORT_CHINESE_BIG5           SORTID = 0x0 // Chinese BIG5 order
	SORT_CHINESE_PRCP           SORTID = 0x0 // PRC Chinese Phonetic order
	SORT_CHINESE_UNICODE        SORTID = 0x1 // Chinese Unicode order (no longer supported)
	SORT_CHINESE_PRC            SORTID = 0x2 // PRC Chinese Stroke Count order
	SORT_CHINESE_BOPOMOFO       SORTID = 0x3 // Traditional Chinese Bopomofo order
	SORT_CHINESE_RADICALSTROKE  SORTID = 0x4 // Traditional Chinese radical/stroke order.
	SORT_KOREAN_KSC             SORTID = 0x0 // Korean KSC order
	SORT_KOREAN_UNICODE         SORTID = 0x1 // Korean Unicode order (no longer supported)
	SORT_GERMAN_PHONE_BOOK      SORTID = 0x1 // German Phone Book order
	SORT_HUNGARIAN_DEFAULT      SORTID = 0x0 // Hungarian Default order
	SORT_HUNGARIAN_TECHNICAL    SORTID = 0x1 // Hungarian Technical order
	SORT_GEORGIAN_TRADITIONAL   SORTID = 0x0 // Georgian Traditional order
	SORT_GEORGIAN_MODERN        SORTID = 0x1 // Georgian Modern order
)

const NLS_VALID_LOCALE_MASK LCID = 0x000fffff

func MAKELANGID(p LANGID_PRIMARY, s LANGID_SUB) LANGID {
	return LANGID(WORD(p) | WORD(s)<<10)
}

func PRIMARYLANGID(lgid LANGID) LANGID_PRIMARY {
	return LANGID_PRIMARY(lgid & 0x3ff)
}

func SUBLANGID(lgid LANGID) LANGID_SUB {
	return LANGID_SUB(lgid >> 10)
}

func MAKELCID(lgid LANGID, srtid SORTID) LCID {
	return LCID((DWORD((srtid)) << 16) | DWORD(lgid))
}

func MAKESORTLCID(lgid LANGID, srtid SORTID, ver WORD) LCID {
	return MAKELCID(lgid, srtid) | (LCID(ver) << 20)
}

func LANGIDFROMLCID(lcid LCID) LANGID {
	return LANGID(lcid)
}

func SORTIDFROMLCID(lcid LCID) SORTID {
	return SORTID(lcid>>16) & 0xf
}

func SORTVERSIONFROMLCID(lcid LCID) WORD {
	return WORD((lcid >> 20) & 0xf)
}
