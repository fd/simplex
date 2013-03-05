package html

type Entity rune

const (
	ENTITY_INVALID  Entity = 0
	ENTITY_Aacute          = 0x00C1
	ENTITY_aacute          = 0x00E1
	ENTITY_Acirc           = 0x00C2
	ENTITY_acirc           = 0x00E2
	ENTITY_acute           = 0x00B4
	ENTITY_aelig           = 0x00E6
	ENTITY_AElig           = 0x00C6
	ENTITY_Agrave          = 0x00C0
	ENTITY_agrave          = 0x00E0
	ENTITY_alefsym         = 0x2135
	ENTITY_Alpha           = 0x0391
	ENTITY_alpha           = 0x03B1
	ENTITY_AMP             = 0x0026
	ENTITY_amp             = 0x0026
	ENTITY_and             = 0x2227
	ENTITY_ang             = 0x2220
	ENTITY_apos            = 0x0027
	ENTITY_Aring           = 0x00C5
	ENTITY_aring           = 0x00E5
	ENTITY_asymp           = 0x2248
	ENTITY_Atilde          = 0x00C3
	ENTITY_atilde          = 0x00E3
	ENTITY_Auml            = 0x00C4
	ENTITY_auml            = 0x00E4
	ENTITY_bdquo           = 0x201E
	ENTITY_Beta            = 0x0392
	ENTITY_beta            = 0x03B2
	ENTITY_brvbar          = 0x00A6
	ENTITY_bull            = 0x2022
	ENTITY_cap             = 0x2229
	ENTITY_Ccedil          = 0x00C7
	ENTITY_ccedil          = 0x00E7
	ENTITY_cedil           = 0x00B8
	ENTITY_cent            = 0x00A2
	ENTITY_Chi             = 0x03A7
	ENTITY_chi             = 0x03C7
	ENTITY_circ            = 0x02C6
	ENTITY_clubs           = 0x2663
	ENTITY_cong            = 0x2245
	ENTITY_COPY            = 0x00A9
	ENTITY_copy            = 0x00A9
	ENTITY_crarr           = 0x21B5
	ENTITY_cup             = 0x222A
	ENTITY_curren          = 0x00A4
	ENTITY_dagger          = 0x2020
	ENTITY_Dagger          = 0x2021
	ENTITY_darr            = 0x2193
	ENTITY_dArr            = 0x21D3
	ENTITY_deg             = 0x00B0
	ENTITY_Delta           = 0x0394
	ENTITY_delta           = 0x03B4
	ENTITY_diams           = 0x2666
	ENTITY_divide          = 0x00F7
	ENTITY_Eacute          = 0x00C9
	ENTITY_eacute          = 0x00E9
	ENTITY_Ecirc           = 0x00CA
	ENTITY_ecirc           = 0x00EA
	ENTITY_Egrave          = 0x00C8
	ENTITY_egrave          = 0x00E8
	ENTITY_empty           = 0x2205
	ENTITY_emsp            = 0x2003
	ENTITY_ensp            = 0x2002
	ENTITY_Epsilon         = 0x0395
	ENTITY_epsilon         = 0x03B5
	ENTITY_equiv           = 0x2261
	ENTITY_Eta             = 0x0397
	ENTITY_eta             = 0x03B7
	ENTITY_ETH             = 0x00D0
	ENTITY_eth             = 0x00F0
	ENTITY_Euml            = 0x00CB
	ENTITY_euml            = 0x00EB
	ENTITY_euro            = 0x20AC
	ENTITY_exist           = 0x2203
	ENTITY_fnof            = 0x0192
	ENTITY_forall          = 0x2200
	ENTITY_frac12          = 0x00BD
	ENTITY_frac14          = 0x00BC
	ENTITY_frac34          = 0x00BE
	ENTITY_frasl           = 0x2044
	ENTITY_Gamma           = 0x0393
	ENTITY_gamma           = 0x03B3
	ENTITY_ge              = 0x2265
	ENTITY_GT              = 0x003E
	ENTITY_gt              = 0x003E
	ENTITY_harr            = 0x2194
	ENTITY_hArr            = 0x21D4
	ENTITY_hearts          = 0x2665
	ENTITY_hellip          = 0x2026
	ENTITY_Iacute          = 0x00CD
	ENTITY_iacute          = 0x00ED
	ENTITY_Icirc           = 0x00CE
	ENTITY_icirc           = 0x00EE
	ENTITY_iexcl           = 0x00A1
	ENTITY_Igrave          = 0x00CC
	ENTITY_igrave          = 0x00EC
	ENTITY_image           = 0x2111
	ENTITY_infin           = 0x221E
	ENTITY_int             = 0x222B
	ENTITY_Iota            = 0x0399
	ENTITY_iota            = 0x03B9
	ENTITY_iquest          = 0x00BF
	ENTITY_isin            = 0x2208
	ENTITY_Iuml            = 0x00CF
	ENTITY_iuml            = 0x00EF
	ENTITY_Kappa           = 0x039A
	ENTITY_kappa           = 0x03BA
	ENTITY_Lambda          = 0x039B
	ENTITY_lambda          = 0x03BB
	ENTITY_lang            = 0x3008
	ENTITY_laquo           = 0x00AB
	ENTITY_larr            = 0x2190
	ENTITY_lArr            = 0x21D0
	ENTITY_lceil           = 0x2308
	ENTITY_ldquo           = 0x201C
	ENTITY_le              = 0x2264
	ENTITY_lfloor          = 0x230A
	ENTITY_lowast          = 0x2217
	ENTITY_loz             = 0x25CA
	ENTITY_lrm             = 0x200E
	ENTITY_lsaquo          = 0x2039
	ENTITY_lsquo           = 0x2018
	ENTITY_LT              = 0x003C
	ENTITY_lt              = 0x003C
	ENTITY_macr            = 0x00AF
	ENTITY_mdash           = 0x2014
	ENTITY_micro           = 0x00B5
	ENTITY_middot          = 0x00B7
	ENTITY_minus           = 0x2212
	ENTITY_Mu              = 0x039C
	ENTITY_mu              = 0x03BC
	ENTITY_nabla           = 0x2207
	ENTITY_nbsp            = 0x00A0
	ENTITY_ndash           = 0x2013
	ENTITY_ne              = 0x2260
	ENTITY_ni              = 0x220B
	ENTITY_not             = 0x00AC
	ENTITY_notin           = 0x2209
	ENTITY_nsub            = 0x2284
	ENTITY_Ntilde          = 0x00D1
	ENTITY_ntilde          = 0x00F1
	ENTITY_Nu              = 0x039D
	ENTITY_nu              = 0x03BD
	ENTITY_Oacute          = 0x00D3
	ENTITY_oacute          = 0x00F3
	ENTITY_Ocirc           = 0x00D4
	ENTITY_ocirc           = 0x00F4
	ENTITY_OElig           = 0x0152
	ENTITY_oelig           = 0x0153
	ENTITY_Ograve          = 0x00D2
	ENTITY_ograve          = 0x00F2
	ENTITY_oline           = 0x203E
	ENTITY_Omega           = 0x03A9
	ENTITY_omega           = 0x03C9
	ENTITY_Omicron         = 0x039F
	ENTITY_omicron         = 0x03BF
	ENTITY_oplus           = 0x2295
	ENTITY_or              = 0x2228
	ENTITY_ordf            = 0x00AA
	ENTITY_ordm            = 0x00BA
	ENTITY_Oslash          = 0x00D8
	ENTITY_oslash          = 0x00F8
	ENTITY_Otilde          = 0x00D5
	ENTITY_otilde          = 0x00F5
	ENTITY_otimes          = 0x2297
	ENTITY_Ouml            = 0x00D6
	ENTITY_ouml            = 0x00F6
	ENTITY_para            = 0x00B6
	ENTITY_part            = 0x2202
	ENTITY_permil          = 0x2030
	ENTITY_perp            = 0x22A5
	ENTITY_Phi             = 0x03A6
	ENTITY_phi             = 0x03C6
	ENTITY_Pi              = 0x03A0
	ENTITY_pi              = 0x03C0
	ENTITY_piv             = 0x03D6
	ENTITY_plusmn          = 0x00B1
	ENTITY_pound           = 0x00A3
	ENTITY_prime           = 0x2032
	ENTITY_Prime           = 0x2033
	ENTITY_prod            = 0x220F
	ENTITY_prop            = 0x221D
	ENTITY_Psi             = 0x03A8
	ENTITY_psi             = 0x03C8
	ENTITY_QUOT            = 0x0022
	ENTITY_quot            = 0x0022
	ENTITY_radic           = 0x221A
	ENTITY_rang            = 0x3009
	ENTITY_raquo           = 0x00BB
	ENTITY_rarr            = 0x2192
	ENTITY_rArr            = 0x21D2
	ENTITY_rceil           = 0x2309
	ENTITY_rdquo           = 0x201D
	ENTITY_real            = 0x211C
	ENTITY_REG             = 0x00AE
	ENTITY_reg             = 0x00AE
	ENTITY_rfloor          = 0x230B
	ENTITY_Rho             = 0x03A1
	ENTITY_rho             = 0x03C1
	ENTITY_rlm             = 0x200F
	ENTITY_rsaquo          = 0x203A
	ENTITY_rsquo           = 0x2019
	ENTITY_sbquo           = 0x201A
	ENTITY_Scaron          = 0x0160
	ENTITY_scaron          = 0x0161
	ENTITY_sdot            = 0x22C5
	ENTITY_sect            = 0x00A7
	ENTITY_shy             = 0x00AD
	ENTITY_Sigma           = 0x03A3
	ENTITY_sigma           = 0x03C3
	ENTITY_sigmaf          = 0x03C2
	ENTITY_sim             = 0x223C
	ENTITY_spades          = 0x2660
	ENTITY_sub             = 0x2282
	ENTITY_sube            = 0x2286
	ENTITY_sum             = 0x2211
	ENTITY_sup             = 0x2283
	ENTITY_sup1            = 0x00B9
	ENTITY_sup2            = 0x00B2
	ENTITY_sup3            = 0x00B3
	ENTITY_supe            = 0x2287
	ENTITY_szlig           = 0x00DF
	ENTITY_Tau             = 0x03A4
	ENTITY_tau             = 0x03C4
	ENTITY_there4          = 0x2234
	ENTITY_Theta           = 0x0398
	ENTITY_theta           = 0x03B8
	ENTITY_thetasym        = 0x03D1
	ENTITY_thinsp          = 0x2009
	ENTITY_THORN           = 0x00DE
	ENTITY_thorn           = 0x00FE
	ENTITY_tilde           = 0x02DC
	ENTITY_times           = 0x00D7
	ENTITY_TRADE           = 0x2122
	ENTITY_trade           = 0x2122
	ENTITY_Uacute          = 0x00DA
	ENTITY_uacute          = 0x00FA
	ENTITY_uarr            = 0x2191
	ENTITY_uArr            = 0x21D1
	ENTITY_Ucirc           = 0x00DB
	ENTITY_ucirc           = 0x00FB
	ENTITY_Ugrave          = 0x00D9
	ENTITY_ugrave          = 0x00F9
	ENTITY_uml             = 0x00A8
	ENTITY_upsih           = 0x03D2
	ENTITY_Upsilon         = 0x03A5
	ENTITY_upsilon         = 0x03C5
	ENTITY_Uuml            = 0x00DC
	ENTITY_uuml            = 0x00FC
	ENTITY_weierp          = 0x2118
	ENTITY_Xi              = 0x039E
	ENTITY_xi              = 0x03BE
	ENTITY_Yacute          = 0x00DD
	ENTITY_yacute          = 0x00FD
	ENTITY_yen             = 0x00A5
	ENTITY_yuml            = 0x00FF
	ENTITY_Yuml            = 0x0178
	ENTITY_Zeta            = 0x0396
	ENTITY_zeta            = 0x03B6
	ENTITY_zwj             = 0x200D
	ENTITY_zwnj            = 0x200C
)

var html_entity_names = map[Entity][]string{
	ENTITY_Aacute:   {"Aacute"},
	ENTITY_aacute:   {"aacute"},
	ENTITY_Acirc:    {"Acirc"},
	ENTITY_acirc:    {"acirc"},
	ENTITY_acute:    {"acute"},
	ENTITY_AElig:    {"AElig"},
	ENTITY_aelig:    {"aelig"},
	ENTITY_Agrave:   {"Agrave"},
	ENTITY_agrave:   {"agrave"},
	ENTITY_alefsym:  {"alefsym"},
	ENTITY_Alpha:    {"Alpha"},
	ENTITY_alpha:    {"alpha"},
	ENTITY_AMP:      {"amp", "AMP"},
	ENTITY_and:      {"and"},
	ENTITY_ang:      {"ang"},
	ENTITY_apos:     {"apos"},
	ENTITY_Aring:    {"Aring"},
	ENTITY_aring:    {"aring"},
	ENTITY_asymp:    {"asymp"},
	ENTITY_Atilde:   {"Atilde"},
	ENTITY_atilde:   {"atilde"},
	ENTITY_Auml:     {"Auml"},
	ENTITY_auml:     {"auml"},
	ENTITY_bdquo:    {"bdquo"},
	ENTITY_Beta:     {"Beta"},
	ENTITY_beta:     {"beta"},
	ENTITY_brvbar:   {"brvbar"},
	ENTITY_bull:     {"bull"},
	ENTITY_cap:      {"cap"},
	ENTITY_Ccedil:   {"Ccedil"},
	ENTITY_ccedil:   {"ccedil"},
	ENTITY_cedil:    {"cedil"},
	ENTITY_cent:     {"cent"},
	ENTITY_Chi:      {"Chi"},
	ENTITY_chi:      {"chi"},
	ENTITY_circ:     {"circ"},
	ENTITY_clubs:    {"clubs"},
	ENTITY_cong:     {"cong"},
	ENTITY_COPY:     {"copy", "COPY"},
	ENTITY_crarr:    {"crarr"},
	ENTITY_cup:      {"cup"},
	ENTITY_curren:   {"curren"},
	ENTITY_Dagger:   {"Dagger"},
	ENTITY_dagger:   {"dagger"},
	ENTITY_dArr:     {"dArr"},
	ENTITY_darr:     {"darr"},
	ENTITY_deg:      {"deg"},
	ENTITY_Delta:    {"Delta"},
	ENTITY_delta:    {"delta"},
	ENTITY_diams:    {"diams"},
	ENTITY_divide:   {"divide"},
	ENTITY_Eacute:   {"Eacute"},
	ENTITY_eacute:   {"eacute"},
	ENTITY_Ecirc:    {"Ecirc"},
	ENTITY_ecirc:    {"ecirc"},
	ENTITY_Egrave:   {"Egrave"},
	ENTITY_egrave:   {"egrave"},
	ENTITY_empty:    {"empty"},
	ENTITY_emsp:     {"emsp"},
	ENTITY_ensp:     {"ensp"},
	ENTITY_Epsilon:  {"Epsilon"},
	ENTITY_epsilon:  {"epsilon"},
	ENTITY_equiv:    {"equiv"},
	ENTITY_Eta:      {"Eta"},
	ENTITY_eta:      {"eta"},
	ENTITY_ETH:      {"ETH"},
	ENTITY_eth:      {"eth"},
	ENTITY_Euml:     {"Euml"},
	ENTITY_euml:     {"euml"},
	ENTITY_euro:     {"euro"},
	ENTITY_exist:    {"exist"},
	ENTITY_fnof:     {"fnof"},
	ENTITY_forall:   {"forall"},
	ENTITY_frac12:   {"frac12"},
	ENTITY_frac14:   {"frac14"},
	ENTITY_frac34:   {"frac34"},
	ENTITY_frasl:    {"frasl"},
	ENTITY_Gamma:    {"Gamma"},
	ENTITY_gamma:    {"gamma"},
	ENTITY_ge:       {"ge"},
	ENTITY_GT:       {"gt", "GT"},
	ENTITY_hArr:     {"hArr"},
	ENTITY_harr:     {"harr"},
	ENTITY_hearts:   {"hearts"},
	ENTITY_hellip:   {"hellip"},
	ENTITY_Iacute:   {"Iacute"},
	ENTITY_iacute:   {"iacute"},
	ENTITY_Icirc:    {"Icirc"},
	ENTITY_icirc:    {"icirc"},
	ENTITY_iexcl:    {"iexcl"},
	ENTITY_Igrave:   {"Igrave"},
	ENTITY_igrave:   {"igrave"},
	ENTITY_image:    {"image"},
	ENTITY_infin:    {"infin"},
	ENTITY_int:      {"int"},
	ENTITY_Iota:     {"Iota"},
	ENTITY_iota:     {"iota"},
	ENTITY_iquest:   {"iquest"},
	ENTITY_isin:     {"isin"},
	ENTITY_Iuml:     {"Iuml"},
	ENTITY_iuml:     {"iuml"},
	ENTITY_Kappa:    {"Kappa"},
	ENTITY_kappa:    {"kappa"},
	ENTITY_Lambda:   {"Lambda"},
	ENTITY_lambda:   {"lambda"},
	ENTITY_lang:     {"lang"},
	ENTITY_laquo:    {"laquo"},
	ENTITY_lArr:     {"lArr"},
	ENTITY_larr:     {"larr"},
	ENTITY_lceil:    {"lceil"},
	ENTITY_ldquo:    {"ldquo"},
	ENTITY_le:       {"le"},
	ENTITY_lfloor:   {"lfloor"},
	ENTITY_lowast:   {"lowast"},
	ENTITY_loz:      {"loz"},
	ENTITY_lrm:      {"lrm"},
	ENTITY_lsaquo:   {"lsaquo"},
	ENTITY_lsquo:    {"lsquo"},
	ENTITY_LT:       {"lt", "LT"},
	ENTITY_macr:     {"macr"},
	ENTITY_mdash:    {"mdash"},
	ENTITY_micro:    {"micro"},
	ENTITY_middot:   {"middot"},
	ENTITY_minus:    {"minus"},
	ENTITY_Mu:       {"Mu"},
	ENTITY_mu:       {"mu"},
	ENTITY_nabla:    {"nabla"},
	ENTITY_nbsp:     {"nbsp"},
	ENTITY_ndash:    {"ndash"},
	ENTITY_ne:       {"ne"},
	ENTITY_ni:       {"ni"},
	ENTITY_not:      {"not"},
	ENTITY_notin:    {"notin"},
	ENTITY_nsub:     {"nsub"},
	ENTITY_Ntilde:   {"Ntilde"},
	ENTITY_ntilde:   {"ntilde"},
	ENTITY_Nu:       {"Nu"},
	ENTITY_nu:       {"nu"},
	ENTITY_Oacute:   {"Oacute"},
	ENTITY_oacute:   {"oacute"},
	ENTITY_Ocirc:    {"Ocirc"},
	ENTITY_ocirc:    {"ocirc"},
	ENTITY_OElig:    {"OElig"},
	ENTITY_oelig:    {"oelig"},
	ENTITY_Ograve:   {"Ograve"},
	ENTITY_ograve:   {"ograve"},
	ENTITY_oline:    {"oline"},
	ENTITY_Omega:    {"Omega"},
	ENTITY_omega:    {"omega"},
	ENTITY_Omicron:  {"Omicron"},
	ENTITY_omicron:  {"omicron"},
	ENTITY_oplus:    {"oplus"},
	ENTITY_or:       {"or"},
	ENTITY_ordf:     {"ordf"},
	ENTITY_ordm:     {"ordm"},
	ENTITY_Oslash:   {"Oslash"},
	ENTITY_oslash:   {"oslash"},
	ENTITY_Otilde:   {"Otilde"},
	ENTITY_otilde:   {"otilde"},
	ENTITY_otimes:   {"otimes"},
	ENTITY_Ouml:     {"Ouml"},
	ENTITY_ouml:     {"ouml"},
	ENTITY_para:     {"para"},
	ENTITY_part:     {"part"},
	ENTITY_permil:   {"permil"},
	ENTITY_perp:     {"perp"},
	ENTITY_Phi:      {"Phi"},
	ENTITY_phi:      {"phi"},
	ENTITY_Pi:       {"Pi"},
	ENTITY_pi:       {"pi"},
	ENTITY_piv:      {"piv"},
	ENTITY_plusmn:   {"plusmn"},
	ENTITY_pound:    {"pound"},
	ENTITY_Prime:    {"Prime"},
	ENTITY_prime:    {"prime"},
	ENTITY_prod:     {"prod"},
	ENTITY_prop:     {"prop"},
	ENTITY_Psi:      {"Psi"},
	ENTITY_psi:      {"psi"},
	ENTITY_QUOT:     {"quot", "QUOT"},
	ENTITY_radic:    {"radic"},
	ENTITY_rang:     {"rang"},
	ENTITY_raquo:    {"raquo"},
	ENTITY_rArr:     {"rArr"},
	ENTITY_rarr:     {"rarr"},
	ENTITY_rceil:    {"rceil"},
	ENTITY_rdquo:    {"rdquo"},
	ENTITY_real:     {"real"},
	ENTITY_REG:      {"reg", "REG"},
	ENTITY_rfloor:   {"rfloor"},
	ENTITY_Rho:      {"Rho"},
	ENTITY_rho:      {"rho"},
	ENTITY_rlm:      {"rlm"},
	ENTITY_rsaquo:   {"rsaquo"},
	ENTITY_rsquo:    {"rsquo"},
	ENTITY_sbquo:    {"sbquo"},
	ENTITY_Scaron:   {"Scaron"},
	ENTITY_scaron:   {"scaron"},
	ENTITY_sdot:     {"sdot"},
	ENTITY_sect:     {"sect"},
	ENTITY_shy:      {"shy"},
	ENTITY_Sigma:    {"Sigma"},
	ENTITY_sigma:    {"sigma"},
	ENTITY_sigmaf:   {"sigmaf"},
	ENTITY_sim:      {"sim"},
	ENTITY_spades:   {"spades"},
	ENTITY_sub:      {"sub"},
	ENTITY_sube:     {"sube"},
	ENTITY_sum:      {"sum"},
	ENTITY_sup1:     {"sup1"},
	ENTITY_sup2:     {"sup2"},
	ENTITY_sup3:     {"sup3"},
	ENTITY_sup:      {"sup"},
	ENTITY_supe:     {"supe"},
	ENTITY_szlig:    {"szlig"},
	ENTITY_Tau:      {"Tau"},
	ENTITY_tau:      {"tau"},
	ENTITY_there4:   {"there4"},
	ENTITY_Theta:    {"Theta"},
	ENTITY_theta:    {"theta"},
	ENTITY_thetasym: {"thetasym"},
	ENTITY_thinsp:   {"thinsp"},
	ENTITY_THORN:    {"THORN"},
	ENTITY_thorn:    {"thorn"},
	ENTITY_tilde:    {"tilde"},
	ENTITY_times:    {"times"},
	ENTITY_TRADE:    {"trade", "TRADE"},
	ENTITY_Uacute:   {"Uacute"},
	ENTITY_uacute:   {"uacute"},
	ENTITY_uArr:     {"uArr"},
	ENTITY_uarr:     {"uarr"},
	ENTITY_Ucirc:    {"Ucirc"},
	ENTITY_ucirc:    {"ucirc"},
	ENTITY_Ugrave:   {"Ugrave"},
	ENTITY_ugrave:   {"ugrave"},
	ENTITY_uml:      {"uml"},
	ENTITY_upsih:    {"upsih"},
	ENTITY_Upsilon:  {"Upsilon"},
	ENTITY_upsilon:  {"upsilon"},
	ENTITY_Uuml:     {"Uuml"},
	ENTITY_uuml:     {"uuml"},
	ENTITY_weierp:   {"weierp"},
	ENTITY_Xi:       {"Xi"},
	ENTITY_xi:       {"xi"},
	ENTITY_Yacute:   {"Yacute"},
	ENTITY_yacute:   {"yacute"},
	ENTITY_yen:      {"yen"},
	ENTITY_Yuml:     {"Yuml"},
	ENTITY_yuml:     {"yuml"},
	ENTITY_Zeta:     {"Zeta"},
	ENTITY_zeta:     {"zeta"},
	ENTITY_zwj:      {"zwj"},
	ENTITY_zwnj:     {"zwnj"},
}

var entity_tokens = func() map[string]Entity {
	m := make(map[string]Entity)
	for e, toks := range html_entity_names {
		for _, tok := range toks {
			m["&"+tok+";"] = e
		}
	}
	return m
}()

func LookupEntity(name string) Entity {
	if e, ok := entity_tokens[name]; ok {
		return e
	}
	return ENTITY_INVALID
}
