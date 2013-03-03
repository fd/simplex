package html

type HTML_ENTITY rune

const (
	HTML_ENTITY_Aacute   HTML_ENTITY = 0x00C1
	HTML_ENTITY_aacute               = 0x00E1
	HTML_ENTITY_Acirc                = 0x00C2
	HTML_ENTITY_acirc                = 0x00E2
	HTML_ENTITY_acute                = 0x00B4
	HTML_ENTITY_aelig                = 0x00E6
	HTML_ENTITY_AElig                = 0x00C6
	HTML_ENTITY_Agrave               = 0x00C0
	HTML_ENTITY_agrave               = 0x00E0
	HTML_ENTITY_alefsym              = 0x2135
	HTML_ENTITY_Alpha                = 0x0391
	HTML_ENTITY_alpha                = 0x03B1
	HTML_ENTITY_AMP                  = 0x0026
	HTML_ENTITY_amp                  = 0x0026
	HTML_ENTITY_and                  = 0x2227
	HTML_ENTITY_ang                  = 0x2220
	HTML_ENTITY_apos                 = 0x0027
	HTML_ENTITY_Aring                = 0x00C5
	HTML_ENTITY_aring                = 0x00E5
	HTML_ENTITY_asymp                = 0x2248
	HTML_ENTITY_Atilde               = 0x00C3
	HTML_ENTITY_atilde               = 0x00E3
	HTML_ENTITY_Auml                 = 0x00C4
	HTML_ENTITY_auml                 = 0x00E4
	HTML_ENTITY_bdquo                = 0x201E
	HTML_ENTITY_Beta                 = 0x0392
	HTML_ENTITY_beta                 = 0x03B2
	HTML_ENTITY_brvbar               = 0x00A6
	HTML_ENTITY_bull                 = 0x2022
	HTML_ENTITY_cap                  = 0x2229
	HTML_ENTITY_Ccedil               = 0x00C7
	HTML_ENTITY_ccedil               = 0x00E7
	HTML_ENTITY_cedil                = 0x00B8
	HTML_ENTITY_cent                 = 0x00A2
	HTML_ENTITY_Chi                  = 0x03A7
	HTML_ENTITY_chi                  = 0x03C7
	HTML_ENTITY_circ                 = 0x02C6
	HTML_ENTITY_clubs                = 0x2663
	HTML_ENTITY_cong                 = 0x2245
	HTML_ENTITY_COPY                 = 0x00A9
	HTML_ENTITY_copy                 = 0x00A9
	HTML_ENTITY_crarr                = 0x21B5
	HTML_ENTITY_cup                  = 0x222A
	HTML_ENTITY_curren               = 0x00A4
	HTML_ENTITY_dagger               = 0x2020
	HTML_ENTITY_Dagger               = 0x2021
	HTML_ENTITY_darr                 = 0x2193
	HTML_ENTITY_dArr                 = 0x21D3
	HTML_ENTITY_deg                  = 0x00B0
	HTML_ENTITY_Delta                = 0x0394
	HTML_ENTITY_delta                = 0x03B4
	HTML_ENTITY_diams                = 0x2666
	HTML_ENTITY_divide               = 0x00F7
	HTML_ENTITY_Eacute               = 0x00C9
	HTML_ENTITY_eacute               = 0x00E9
	HTML_ENTITY_Ecirc                = 0x00CA
	HTML_ENTITY_ecirc                = 0x00EA
	HTML_ENTITY_Egrave               = 0x00C8
	HTML_ENTITY_egrave               = 0x00E8
	HTML_ENTITY_empty                = 0x2205
	HTML_ENTITY_emsp                 = 0x2003
	HTML_ENTITY_ensp                 = 0x2002
	HTML_ENTITY_Epsilon              = 0x0395
	HTML_ENTITY_epsilon              = 0x03B5
	HTML_ENTITY_equiv                = 0x2261
	HTML_ENTITY_Eta                  = 0x0397
	HTML_ENTITY_eta                  = 0x03B7
	HTML_ENTITY_ETH                  = 0x00D0
	HTML_ENTITY_eth                  = 0x00F0
	HTML_ENTITY_Euml                 = 0x00CB
	HTML_ENTITY_euml                 = 0x00EB
	HTML_ENTITY_euro                 = 0x20AC
	HTML_ENTITY_exist                = 0x2203
	HTML_ENTITY_fnof                 = 0x0192
	HTML_ENTITY_forall               = 0x2200
	HTML_ENTITY_frac12               = 0x00BD
	HTML_ENTITY_frac14               = 0x00BC
	HTML_ENTITY_frac34               = 0x00BE
	HTML_ENTITY_frasl                = 0x2044
	HTML_ENTITY_Gamma                = 0x0393
	HTML_ENTITY_gamma                = 0x03B3
	HTML_ENTITY_ge                   = 0x2265
	HTML_ENTITY_GT                   = 0x003E
	HTML_ENTITY_gt                   = 0x003E
	HTML_ENTITY_harr                 = 0x2194
	HTML_ENTITY_hArr                 = 0x21D4
	HTML_ENTITY_hearts               = 0x2665
	HTML_ENTITY_hellip               = 0x2026
	HTML_ENTITY_Iacute               = 0x00CD
	HTML_ENTITY_iacute               = 0x00ED
	HTML_ENTITY_Icirc                = 0x00CE
	HTML_ENTITY_icirc                = 0x00EE
	HTML_ENTITY_iexcl                = 0x00A1
	HTML_ENTITY_Igrave               = 0x00CC
	HTML_ENTITY_igrave               = 0x00EC
	HTML_ENTITY_image                = 0x2111
	HTML_ENTITY_infin                = 0x221E
	HTML_ENTITY_int                  = 0x222B
	HTML_ENTITY_Iota                 = 0x0399
	HTML_ENTITY_iota                 = 0x03B9
	HTML_ENTITY_iquest               = 0x00BF
	HTML_ENTITY_isin                 = 0x2208
	HTML_ENTITY_Iuml                 = 0x00CF
	HTML_ENTITY_iuml                 = 0x00EF
	HTML_ENTITY_Kappa                = 0x039A
	HTML_ENTITY_kappa                = 0x03BA
	HTML_ENTITY_Lambda               = 0x039B
	HTML_ENTITY_lambda               = 0x03BB
	HTML_ENTITY_lang                 = 0x3008
	HTML_ENTITY_laquo                = 0x00AB
	HTML_ENTITY_larr                 = 0x2190
	HTML_ENTITY_lArr                 = 0x21D0
	HTML_ENTITY_lceil                = 0x2308
	HTML_ENTITY_ldquo                = 0x201C
	HTML_ENTITY_le                   = 0x2264
	HTML_ENTITY_lfloor               = 0x230A
	HTML_ENTITY_lowast               = 0x2217
	HTML_ENTITY_loz                  = 0x25CA
	HTML_ENTITY_lrm                  = 0x200E
	HTML_ENTITY_lsaquo               = 0x2039
	HTML_ENTITY_lsquo                = 0x2018
	HTML_ENTITY_LT                   = 0x003C
	HTML_ENTITY_lt                   = 0x003C
	HTML_ENTITY_macr                 = 0x00AF
	HTML_ENTITY_mdash                = 0x2014
	HTML_ENTITY_micro                = 0x00B5
	HTML_ENTITY_middot               = 0x00B7
	HTML_ENTITY_minus                = 0x2212
	HTML_ENTITY_Mu                   = 0x039C
	HTML_ENTITY_mu                   = 0x03BC
	HTML_ENTITY_nabla                = 0x2207
	HTML_ENTITY_nbsp                 = 0x00A0
	HTML_ENTITY_ndash                = 0x2013
	HTML_ENTITY_ne                   = 0x2260
	HTML_ENTITY_ni                   = 0x220B
	HTML_ENTITY_not                  = 0x00AC
	HTML_ENTITY_notin                = 0x2209
	HTML_ENTITY_nsub                 = 0x2284
	HTML_ENTITY_Ntilde               = 0x00D1
	HTML_ENTITY_ntilde               = 0x00F1
	HTML_ENTITY_Nu                   = 0x039D
	HTML_ENTITY_nu                   = 0x03BD
	HTML_ENTITY_Oacute               = 0x00D3
	HTML_ENTITY_oacute               = 0x00F3
	HTML_ENTITY_Ocirc                = 0x00D4
	HTML_ENTITY_ocirc                = 0x00F4
	HTML_ENTITY_OElig                = 0x0152
	HTML_ENTITY_oelig                = 0x0153
	HTML_ENTITY_Ograve               = 0x00D2
	HTML_ENTITY_ograve               = 0x00F2
	HTML_ENTITY_oline                = 0x203E
	HTML_ENTITY_Omega                = 0x03A9
	HTML_ENTITY_omega                = 0x03C9
	HTML_ENTITY_Omicron              = 0x039F
	HTML_ENTITY_omicron              = 0x03BF
	HTML_ENTITY_oplus                = 0x2295
	HTML_ENTITY_or                   = 0x2228
	HTML_ENTITY_ordf                 = 0x00AA
	HTML_ENTITY_ordm                 = 0x00BA
	HTML_ENTITY_Oslash               = 0x00D8
	HTML_ENTITY_oslash               = 0x00F8
	HTML_ENTITY_Otilde               = 0x00D5
	HTML_ENTITY_otilde               = 0x00F5
	HTML_ENTITY_otimes               = 0x2297
	HTML_ENTITY_Ouml                 = 0x00D6
	HTML_ENTITY_ouml                 = 0x00F6
	HTML_ENTITY_para                 = 0x00B6
	HTML_ENTITY_part                 = 0x2202
	HTML_ENTITY_permil               = 0x2030
	HTML_ENTITY_perp                 = 0x22A5
	HTML_ENTITY_Phi                  = 0x03A6
	HTML_ENTITY_phi                  = 0x03C6
	HTML_ENTITY_Pi                   = 0x03A0
	HTML_ENTITY_pi                   = 0x03C0
	HTML_ENTITY_piv                  = 0x03D6
	HTML_ENTITY_plusmn               = 0x00B1
	HTML_ENTITY_pound                = 0x00A3
	HTML_ENTITY_prime                = 0x2032
	HTML_ENTITY_Prime                = 0x2033
	HTML_ENTITY_prod                 = 0x220F
	HTML_ENTITY_prop                 = 0x221D
	HTML_ENTITY_Psi                  = 0x03A8
	HTML_ENTITY_psi                  = 0x03C8
	HTML_ENTITY_QUOT                 = 0x0022
	HTML_ENTITY_quot                 = 0x0022
	HTML_ENTITY_radic                = 0x221A
	HTML_ENTITY_rang                 = 0x3009
	HTML_ENTITY_raquo                = 0x00BB
	HTML_ENTITY_rarr                 = 0x2192
	HTML_ENTITY_rArr                 = 0x21D2
	HTML_ENTITY_rceil                = 0x2309
	HTML_ENTITY_rdquo                = 0x201D
	HTML_ENTITY_real                 = 0x211C
	HTML_ENTITY_REG                  = 0x00AE
	HTML_ENTITY_reg                  = 0x00AE
	HTML_ENTITY_rfloor               = 0x230B
	HTML_ENTITY_Rho                  = 0x03A1
	HTML_ENTITY_rho                  = 0x03C1
	HTML_ENTITY_rlm                  = 0x200F
	HTML_ENTITY_rsaquo               = 0x203A
	HTML_ENTITY_rsquo                = 0x2019
	HTML_ENTITY_sbquo                = 0x201A
	HTML_ENTITY_Scaron               = 0x0160
	HTML_ENTITY_scaron               = 0x0161
	HTML_ENTITY_sdot                 = 0x22C5
	HTML_ENTITY_sect                 = 0x00A7
	HTML_ENTITY_shy                  = 0x00AD
	HTML_ENTITY_Sigma                = 0x03A3
	HTML_ENTITY_sigma                = 0x03C3
	HTML_ENTITY_sigmaf               = 0x03C2
	HTML_ENTITY_sim                  = 0x223C
	HTML_ENTITY_spades               = 0x2660
	HTML_ENTITY_sub                  = 0x2282
	HTML_ENTITY_sube                 = 0x2286
	HTML_ENTITY_sum                  = 0x2211
	HTML_ENTITY_sup                  = 0x2283
	HTML_ENTITY_sup1                 = 0x00B9
	HTML_ENTITY_sup2                 = 0x00B2
	HTML_ENTITY_sup3                 = 0x00B3
	HTML_ENTITY_supe                 = 0x2287
	HTML_ENTITY_szlig                = 0x00DF
	HTML_ENTITY_Tau                  = 0x03A4
	HTML_ENTITY_tau                  = 0x03C4
	HTML_ENTITY_there4               = 0x2234
	HTML_ENTITY_Theta                = 0x0398
	HTML_ENTITY_theta                = 0x03B8
	HTML_ENTITY_thetasym             = 0x03D1
	HTML_ENTITY_thinsp               = 0x2009
	HTML_ENTITY_THORN                = 0x00DE
	HTML_ENTITY_thorn                = 0x00FE
	HTML_ENTITY_tilde                = 0x02DC
	HTML_ENTITY_times                = 0x00D7
	HTML_ENTITY_TRADE                = 0x2122
	HTML_ENTITY_trade                = 0x2122
	HTML_ENTITY_Uacute               = 0x00DA
	HTML_ENTITY_uacute               = 0x00FA
	HTML_ENTITY_uarr                 = 0x2191
	HTML_ENTITY_uArr                 = 0x21D1
	HTML_ENTITY_Ucirc                = 0x00DB
	HTML_ENTITY_ucirc                = 0x00FB
	HTML_ENTITY_Ugrave               = 0x00D9
	HTML_ENTITY_ugrave               = 0x00F9
	HTML_ENTITY_uml                  = 0x00A8
	HTML_ENTITY_upsih                = 0x03D2
	HTML_ENTITY_Upsilon              = 0x03A5
	HTML_ENTITY_upsilon              = 0x03C5
	HTML_ENTITY_Uuml                 = 0x00DC
	HTML_ENTITY_uuml                 = 0x00FC
	HTML_ENTITY_weierp               = 0x2118
	HTML_ENTITY_Xi                   = 0x039E
	HTML_ENTITY_xi                   = 0x03BE
	HTML_ENTITY_Yacute               = 0x00DD
	HTML_ENTITY_yacute               = 0x00FD
	HTML_ENTITY_yen                  = 0x00A5
	HTML_ENTITY_yuml                 = 0x00FF
	HTML_ENTITY_Yuml                 = 0x0178
	HTML_ENTITY_Zeta                 = 0x0396
	HTML_ENTITY_zeta                 = 0x03B6
	HTML_ENTITY_zwj                  = 0x200D
	HTML_ENTITY_zwnj                 = 0x200C
)

var html_entity_names = map[HTML_ENTITY][]string{
	HTML_ENTITY_Aacute:   {"Aacute"},
	HTML_ENTITY_aacute:   {"aacute"},
	HTML_ENTITY_Acirc:    {"Acirc"},
	HTML_ENTITY_acirc:    {"acirc"},
	HTML_ENTITY_acute:    {"acute"},
	HTML_ENTITY_AElig:    {"AElig"},
	HTML_ENTITY_aelig:    {"aelig"},
	HTML_ENTITY_Agrave:   {"Agrave"},
	HTML_ENTITY_agrave:   {"agrave"},
	HTML_ENTITY_alefsym:  {"alefsym"},
	HTML_ENTITY_Alpha:    {"Alpha"},
	HTML_ENTITY_alpha:    {"alpha"},
	HTML_ENTITY_AMP:      {"amp", "AMP"},
	HTML_ENTITY_and:      {"and"},
	HTML_ENTITY_ang:      {"ang"},
	HTML_ENTITY_apos:     {"apos"},
	HTML_ENTITY_Aring:    {"Aring"},
	HTML_ENTITY_aring:    {"aring"},
	HTML_ENTITY_asymp:    {"asymp"},
	HTML_ENTITY_Atilde:   {"Atilde"},
	HTML_ENTITY_atilde:   {"atilde"},
	HTML_ENTITY_Auml:     {"Auml"},
	HTML_ENTITY_auml:     {"auml"},
	HTML_ENTITY_bdquo:    {"bdquo"},
	HTML_ENTITY_Beta:     {"Beta"},
	HTML_ENTITY_beta:     {"beta"},
	HTML_ENTITY_brvbar:   {"brvbar"},
	HTML_ENTITY_bull:     {"bull"},
	HTML_ENTITY_cap:      {"cap"},
	HTML_ENTITY_Ccedil:   {"Ccedil"},
	HTML_ENTITY_ccedil:   {"ccedil"},
	HTML_ENTITY_cedil:    {"cedil"},
	HTML_ENTITY_cent:     {"cent"},
	HTML_ENTITY_Chi:      {"Chi"},
	HTML_ENTITY_chi:      {"chi"},
	HTML_ENTITY_circ:     {"circ"},
	HTML_ENTITY_clubs:    {"clubs"},
	HTML_ENTITY_cong:     {"cong"},
	HTML_ENTITY_COPY:     {"copy", "COPY"},
	HTML_ENTITY_crarr:    {"crarr"},
	HTML_ENTITY_cup:      {"cup"},
	HTML_ENTITY_curren:   {"curren"},
	HTML_ENTITY_Dagger:   {"Dagger"},
	HTML_ENTITY_dagger:   {"dagger"},
	HTML_ENTITY_dArr:     {"dArr"},
	HTML_ENTITY_darr:     {"darr"},
	HTML_ENTITY_deg:      {"deg"},
	HTML_ENTITY_Delta:    {"Delta"},
	HTML_ENTITY_delta:    {"delta"},
	HTML_ENTITY_diams:    {"diams"},
	HTML_ENTITY_divide:   {"divide"},
	HTML_ENTITY_Eacute:   {"Eacute"},
	HTML_ENTITY_eacute:   {"eacute"},
	HTML_ENTITY_Ecirc:    {"Ecirc"},
	HTML_ENTITY_ecirc:    {"ecirc"},
	HTML_ENTITY_Egrave:   {"Egrave"},
	HTML_ENTITY_egrave:   {"egrave"},
	HTML_ENTITY_empty:    {"empty"},
	HTML_ENTITY_emsp:     {"emsp"},
	HTML_ENTITY_ensp:     {"ensp"},
	HTML_ENTITY_Epsilon:  {"Epsilon"},
	HTML_ENTITY_epsilon:  {"epsilon"},
	HTML_ENTITY_equiv:    {"equiv"},
	HTML_ENTITY_Eta:      {"Eta"},
	HTML_ENTITY_eta:      {"eta"},
	HTML_ENTITY_ETH:      {"ETH"},
	HTML_ENTITY_eth:      {"eth"},
	HTML_ENTITY_Euml:     {"Euml"},
	HTML_ENTITY_euml:     {"euml"},
	HTML_ENTITY_euro:     {"euro"},
	HTML_ENTITY_exist:    {"exist"},
	HTML_ENTITY_fnof:     {"fnof"},
	HTML_ENTITY_forall:   {"forall"},
	HTML_ENTITY_frac12:   {"frac12"},
	HTML_ENTITY_frac14:   {"frac14"},
	HTML_ENTITY_frac34:   {"frac34"},
	HTML_ENTITY_frasl:    {"frasl"},
	HTML_ENTITY_Gamma:    {"Gamma"},
	HTML_ENTITY_gamma:    {"gamma"},
	HTML_ENTITY_ge:       {"ge"},
	HTML_ENTITY_GT:       {"gt", "GT"},
	HTML_ENTITY_hArr:     {"hArr"},
	HTML_ENTITY_harr:     {"harr"},
	HTML_ENTITY_hearts:   {"hearts"},
	HTML_ENTITY_hellip:   {"hellip"},
	HTML_ENTITY_Iacute:   {"Iacute"},
	HTML_ENTITY_iacute:   {"iacute"},
	HTML_ENTITY_Icirc:    {"Icirc"},
	HTML_ENTITY_icirc:    {"icirc"},
	HTML_ENTITY_iexcl:    {"iexcl"},
	HTML_ENTITY_Igrave:   {"Igrave"},
	HTML_ENTITY_igrave:   {"igrave"},
	HTML_ENTITY_image:    {"image"},
	HTML_ENTITY_infin:    {"infin"},
	HTML_ENTITY_int:      {"int"},
	HTML_ENTITY_Iota:     {"Iota"},
	HTML_ENTITY_iota:     {"iota"},
	HTML_ENTITY_iquest:   {"iquest"},
	HTML_ENTITY_isin:     {"isin"},
	HTML_ENTITY_Iuml:     {"Iuml"},
	HTML_ENTITY_iuml:     {"iuml"},
	HTML_ENTITY_Kappa:    {"Kappa"},
	HTML_ENTITY_kappa:    {"kappa"},
	HTML_ENTITY_Lambda:   {"Lambda"},
	HTML_ENTITY_lambda:   {"lambda"},
	HTML_ENTITY_lang:     {"lang"},
	HTML_ENTITY_laquo:    {"laquo"},
	HTML_ENTITY_lArr:     {"lArr"},
	HTML_ENTITY_larr:     {"larr"},
	HTML_ENTITY_lceil:    {"lceil"},
	HTML_ENTITY_ldquo:    {"ldquo"},
	HTML_ENTITY_le:       {"le"},
	HTML_ENTITY_lfloor:   {"lfloor"},
	HTML_ENTITY_lowast:   {"lowast"},
	HTML_ENTITY_loz:      {"loz"},
	HTML_ENTITY_lrm:      {"lrm"},
	HTML_ENTITY_lsaquo:   {"lsaquo"},
	HTML_ENTITY_lsquo:    {"lsquo"},
	HTML_ENTITY_LT:       {"lt", "LT"},
	HTML_ENTITY_macr:     {"macr"},
	HTML_ENTITY_mdash:    {"mdash"},
	HTML_ENTITY_micro:    {"micro"},
	HTML_ENTITY_middot:   {"middot"},
	HTML_ENTITY_minus:    {"minus"},
	HTML_ENTITY_Mu:       {"Mu"},
	HTML_ENTITY_mu:       {"mu"},
	HTML_ENTITY_nabla:    {"nabla"},
	HTML_ENTITY_nbsp:     {"nbsp"},
	HTML_ENTITY_ndash:    {"ndash"},
	HTML_ENTITY_ne:       {"ne"},
	HTML_ENTITY_ni:       {"ni"},
	HTML_ENTITY_not:      {"not"},
	HTML_ENTITY_notin:    {"notin"},
	HTML_ENTITY_nsub:     {"nsub"},
	HTML_ENTITY_Ntilde:   {"Ntilde"},
	HTML_ENTITY_ntilde:   {"ntilde"},
	HTML_ENTITY_Nu:       {"Nu"},
	HTML_ENTITY_nu:       {"nu"},
	HTML_ENTITY_Oacute:   {"Oacute"},
	HTML_ENTITY_oacute:   {"oacute"},
	HTML_ENTITY_Ocirc:    {"Ocirc"},
	HTML_ENTITY_ocirc:    {"ocirc"},
	HTML_ENTITY_OElig:    {"OElig"},
	HTML_ENTITY_oelig:    {"oelig"},
	HTML_ENTITY_Ograve:   {"Ograve"},
	HTML_ENTITY_ograve:   {"ograve"},
	HTML_ENTITY_oline:    {"oline"},
	HTML_ENTITY_Omega:    {"Omega"},
	HTML_ENTITY_omega:    {"omega"},
	HTML_ENTITY_Omicron:  {"Omicron"},
	HTML_ENTITY_omicron:  {"omicron"},
	HTML_ENTITY_oplus:    {"oplus"},
	HTML_ENTITY_or:       {"or"},
	HTML_ENTITY_ordf:     {"ordf"},
	HTML_ENTITY_ordm:     {"ordm"},
	HTML_ENTITY_Oslash:   {"Oslash"},
	HTML_ENTITY_oslash:   {"oslash"},
	HTML_ENTITY_Otilde:   {"Otilde"},
	HTML_ENTITY_otilde:   {"otilde"},
	HTML_ENTITY_otimes:   {"otimes"},
	HTML_ENTITY_Ouml:     {"Ouml"},
	HTML_ENTITY_ouml:     {"ouml"},
	HTML_ENTITY_para:     {"para"},
	HTML_ENTITY_part:     {"part"},
	HTML_ENTITY_permil:   {"permil"},
	HTML_ENTITY_perp:     {"perp"},
	HTML_ENTITY_Phi:      {"Phi"},
	HTML_ENTITY_phi:      {"phi"},
	HTML_ENTITY_Pi:       {"Pi"},
	HTML_ENTITY_pi:       {"pi"},
	HTML_ENTITY_piv:      {"piv"},
	HTML_ENTITY_plusmn:   {"plusmn"},
	HTML_ENTITY_pound:    {"pound"},
	HTML_ENTITY_Prime:    {"Prime"},
	HTML_ENTITY_prime:    {"prime"},
	HTML_ENTITY_prod:     {"prod"},
	HTML_ENTITY_prop:     {"prop"},
	HTML_ENTITY_Psi:      {"Psi"},
	HTML_ENTITY_psi:      {"psi"},
	HTML_ENTITY_QUOT:     {"quot", "QUOT"},
	HTML_ENTITY_radic:    {"radic"},
	HTML_ENTITY_rang:     {"rang"},
	HTML_ENTITY_raquo:    {"raquo"},
	HTML_ENTITY_rArr:     {"rArr"},
	HTML_ENTITY_rarr:     {"rarr"},
	HTML_ENTITY_rceil:    {"rceil"},
	HTML_ENTITY_rdquo:    {"rdquo"},
	HTML_ENTITY_real:     {"real"},
	HTML_ENTITY_REG:      {"reg", "REG"},
	HTML_ENTITY_rfloor:   {"rfloor"},
	HTML_ENTITY_Rho:      {"Rho"},
	HTML_ENTITY_rho:      {"rho"},
	HTML_ENTITY_rlm:      {"rlm"},
	HTML_ENTITY_rsaquo:   {"rsaquo"},
	HTML_ENTITY_rsquo:    {"rsquo"},
	HTML_ENTITY_sbquo:    {"sbquo"},
	HTML_ENTITY_Scaron:   {"Scaron"},
	HTML_ENTITY_scaron:   {"scaron"},
	HTML_ENTITY_sdot:     {"sdot"},
	HTML_ENTITY_sect:     {"sect"},
	HTML_ENTITY_shy:      {"shy"},
	HTML_ENTITY_Sigma:    {"Sigma"},
	HTML_ENTITY_sigma:    {"sigma"},
	HTML_ENTITY_sigmaf:   {"sigmaf"},
	HTML_ENTITY_sim:      {"sim"},
	HTML_ENTITY_spades:   {"spades"},
	HTML_ENTITY_sub:      {"sub"},
	HTML_ENTITY_sube:     {"sube"},
	HTML_ENTITY_sum:      {"sum"},
	HTML_ENTITY_sup1:     {"sup1"},
	HTML_ENTITY_sup2:     {"sup2"},
	HTML_ENTITY_sup3:     {"sup3"},
	HTML_ENTITY_sup:      {"sup"},
	HTML_ENTITY_supe:     {"supe"},
	HTML_ENTITY_szlig:    {"szlig"},
	HTML_ENTITY_Tau:      {"Tau"},
	HTML_ENTITY_tau:      {"tau"},
	HTML_ENTITY_there4:   {"there4"},
	HTML_ENTITY_Theta:    {"Theta"},
	HTML_ENTITY_theta:    {"theta"},
	HTML_ENTITY_thetasym: {"thetasym"},
	HTML_ENTITY_thinsp:   {"thinsp"},
	HTML_ENTITY_THORN:    {"THORN"},
	HTML_ENTITY_thorn:    {"thorn"},
	HTML_ENTITY_tilde:    {"tilde"},
	HTML_ENTITY_times:    {"times"},
	HTML_ENTITY_TRADE:    {"trade", "TRADE"},
	HTML_ENTITY_Uacute:   {"Uacute"},
	HTML_ENTITY_uacute:   {"uacute"},
	HTML_ENTITY_uArr:     {"uArr"},
	HTML_ENTITY_uarr:     {"uarr"},
	HTML_ENTITY_Ucirc:    {"Ucirc"},
	HTML_ENTITY_ucirc:    {"ucirc"},
	HTML_ENTITY_Ugrave:   {"Ugrave"},
	HTML_ENTITY_ugrave:   {"ugrave"},
	HTML_ENTITY_uml:      {"uml"},
	HTML_ENTITY_upsih:    {"upsih"},
	HTML_ENTITY_Upsilon:  {"Upsilon"},
	HTML_ENTITY_upsilon:  {"upsilon"},
	HTML_ENTITY_Uuml:     {"Uuml"},
	HTML_ENTITY_uuml:     {"uuml"},
	HTML_ENTITY_weierp:   {"weierp"},
	HTML_ENTITY_Xi:       {"Xi"},
	HTML_ENTITY_xi:       {"xi"},
	HTML_ENTITY_Yacute:   {"Yacute"},
	HTML_ENTITY_yacute:   {"yacute"},
	HTML_ENTITY_yen:      {"yen"},
	HTML_ENTITY_Yuml:     {"Yuml"},
	HTML_ENTITY_yuml:     {"yuml"},
	HTML_ENTITY_Zeta:     {"Zeta"},
	HTML_ENTITY_zeta:     {"zeta"},
	HTML_ENTITY_zwj:      {"zwj"},
	HTML_ENTITY_zwnj:     {"zwnj"},
}
