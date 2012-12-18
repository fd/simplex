package views

import (
	import_2 "github.com/fd/w/data"
	import_1 "github.com/fd/w/example/apps/orakel/helpers"
	import_3 "github.com/fd/w/runtime"
)

func Footer(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   1
	   column: 1
	*/
	buf.Write(import_3.HTML("  <footer>\n\n  <div data-role=\"row\" class=\"row keypoints\">\n\n    <div class=\"wrapper keypoints\">\n     "))
	buf.Write(import_3.HTML(" <div data-role=\"column\" class=\"quality\">\n        <h4>Service</h4>\n        <p style=\"text-align: cen"))
	buf.Write(import_3.HTML("ter;\" align=\"center\">We treat our clients as we would like to be treated ourselves.. a personable, k"))
	buf.Write(import_3.HTML("nowledgeable team offering quality products, sensibly priced with fast and efficient service.</p>\n  "))
	buf.Write(import_3.HTML("    </div>\n      <div data-role=\"column\" class=\"experience\">\n        <h4>Experience</h4>\n        <p "))
	buf.Write(import_3.HTML("style=\"text-align: center;\" align=\"center\">Founded in&nbsp;1996, Orakel launched the first company i"))
	buf.Write(import_3.HTML("n Europe&nbsp;to&nbsp;offer&nbsp;&nbsp;a unique range of products.</p>\n<p style=\"text-align: center;"))
	buf.Write(import_3.HTML("\" align=\"center\">By designing and manufacturing our products in-house,&nbsp;our team, well respected"))
	buf.Write(import_3.HTML(" within the industry for their level of expertise&nbsp;is able to respond swiftly to&nbsp;the needs "))
	buf.Write(import_3.HTML("of our clients.&nbsp;</p>\n      </div>\n      <div data-role=\"column\" class=\"supply\">\n        <h4>Inn"))
	buf.Write(import_3.HTML("ovation</h4>\n        <p style=\"text-align: center;\" align=\"center\">We believe Orakel's succes lies i"))
	buf.Write(import_3.HTML("n meeting the needs of our clients with an constantly evolving product range, every new situation wh"))
	buf.Write(import_3.HTML("ich presents itself fuels fresh inspirational answers.<br /> <br /> You can follow our newest develo"))
	buf.Write(import_3.HTML("pments through our website and facebook.</p>\n      </div>\n      <div data-role=\"column\" class=\"last "))
	buf.Write(import_3.HTML("sustainability\">\n        <h4>Sustainability</h4>\n        <p style=\"text-align: center;\" align=\"cente"))
	buf.Write(import_3.HTML("r\">At Orakel we work to support the ethos of &nbsp;'People, Planet, Profit'. We balance the care of "))
	buf.Write(import_3.HTML("our clients, the welfare of our staff, working with respect for the earth's resources and reducing o"))
	buf.Write(import_3.HTML("ur carbon footprint whilst contributing to a&nbsp; flourishing economy.</p>\n      </div>\n    </div>\n"))
	buf.Write(import_3.HTML("\n    <div class=\"wrapper partners\">\n      <ul>\n        <li class=\"intro\">our partners</li>\n        <"))
	buf.Write(import_3.HTML("li class=\"fit\"><a href=\"http://www.flandersinvestmentandtrade.com\" target=\"_blank\">Flanders Investme"))
	buf.Write(import_3.HTML("nt and Trade</a></li>\n        <li class=\"besa\"><a href=\"http://www.b-esa.be\" target=\"_blank\">BESA</a"))
	buf.Write(import_3.HTML("></li>\n      </ul>\n    </div>\n\n  </div>\n\n  <div data-role=\"row\" class=\"row contact\">\n\n    <div class"))
	buf.Write(import_3.HTML("=\"wrapper\">\n      <div data-role=\"column\" class=\"column one\">\n        <h4>Address</h4>\n        <p>Ve"))
	buf.Write(import_3.HTML("ldenstraat 14bis<br />2470 Retie<br />Belgium</p>\n      </div>\n      <div data-role=\"column\" class=\""))
	buf.Write(import_3.HTML("column two\">\n        <p>t +32 (0) 14 38 80 80 <br />f +32 (0) 14 38 80 10 <br /><a href=\"mailto:orak"))
	buf.Write(import_3.HTML("el@orakel.com\">orakel@orakel.com</a></p>\n      </div>\n      <div data-role=\"column\" class=\"column th"))
	buf.Write(import_3.HTML("ree\">\n\n      </div>\n      <div data-role=\"column\" class=\"column four\">\n        <h4>At your service</"))
	buf.Write(import_3.HTML("h4>\n        <p>You can find us in many European countries. Please click on your country.</p>\n       "))
	buf.Write(import_3.HTML(" <nav>\n\n          <a href=\"http://be.orakel.com\">Belgium</a>\n\n          <a href=\"http://orakel.com\">"))
	buf.Write(import_3.HTML("Europe</a>\n\n          <a href=\"http://fr.orakel.com\">France</a>\n\n          <a href=\"http://de.orakel"))
	buf.Write(import_3.HTML(".com\">Germany</a>\n\n          <a href=\"http://gr.orakel.com\">Greece</a>\n\n          <a href=\"http://hu"))
	buf.Write(import_3.HTML(".orakel.com\">Hungary</a>\n\n          <a href=\"http://it.orakel.com\">Italy</a>\n\n          <a href=\"htt"))
	buf.Write(import_3.HTML("p://nl.orakel.com\">Netherlands</a>\n\n          <a href=\"http://no.orakel.com\">Norway</a>\n\n          <"))
	buf.Write(import_3.HTML("a href=\"http://pl.orakel.com\">Poland</a>\n\n          <a href=\"http://pt.orakel.com\">Portugal</a>\n\n   "))
	buf.Write(import_3.HTML("       <a href=\"http://es.orakel.com\">Spain</a>\n\n          <a href=\"http://tr.orakel.com\">Turkey</a>"))
	buf.Write(import_3.HTML("\n\n          <a href=\"http://uk.orakel.com\">United Kingdom</a>\n\n        </nav>\n      </div>\n    </div"))
	buf.Write(import_3.HTML(">\n\n  </div>\n\n  <div data-role=\"row\" class=\"row credits\">\n\n    <div class=\"wrapper\">\n      <div data-"))
	buf.Write(import_3.HTML("role=\"column\" class=\"column left\">\n        copyright Orakel -\n\n          <a href=\"/disclaimer\">discl"))
	buf.Write(import_3.HTML("aimer</a> /\n          <a href=\"/terms-and-conditions\">terms and conditions</a> /\n          <a href=\""))
	buf.Write(import_3.HTML("/sitemap\">sitemap</a>\n\n      </div>\n      <div data-role=\"column\" class=\"column right\">\n        webs"))
	buf.Write(import_3.HTML("ite by\n        <a href=\"http://mrhenry.be\" target=\"_blank\">Mr. Henry</a>\n      </div>\n    </div>\n  <"))
	buf.Write(import_3.HTML("/div>\n\n</footer>\n\n\n<script src=\"/javascripts/javascript.js\"></script>\n\n<!--[if lt IE 7 ]>\n  <script "))
	buf.Write(import_3.HTML("src=\"//ajax.googleapis.com/ajax/libs/chrome-frame/1.0.2/CFInstall.min.js\"></script>\n  <script>window"))
	buf.Write(import_3.HTML(".attachEvent(\"onload\",function(){CFInstall.check({mode:\"overlay\"})})</script>\n<![endif]-->\n\n<script "))
	buf.Write(import_3.HTML("type=\"text/javascript\">\n  var _gaq = _gaq || [];\n  _gaq.push(['_setAccount', 'UA-36616061-1']);\n  _g"))
	buf.Write(import_3.HTML("aq.push(['_trackPageview']);\n  (function() {\n    var ga = document.createElement('script'); ga.type "))
	buf.Write(import_3.HTML("= 'text/javascript'; ga.async = true;\n    ga.src = ('https:' == document.location.protocol ? 'https:"))
	buf.Write(import_3.HTML("//ssl' : 'http://www') + '.google-analytics.com/ga.js';\n    var s = document.getElementsByTagName('s"))
	buf.Write(import_3.HTML("cript')[0]; s.parentNode.insertBefore(ga, s);\n  })();\n</script>\n<script type=\"text/javascript\">\n  va"))
	buf.Write(import_3.HTML("r _gaq = _gaq || [];\n  _gaq.push(['_setAccount', 'UA-123802-2']);\n  _gaq.push(['_trackPageview']);\n "))
	buf.Write(import_3.HTML(" (function() {\n    var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async "))
	buf.Write(import_3.HTML("= true;\n    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.goo"))
	buf.Write(import_3.HTML("gle-analytics.com/ga.js';\n    var s = document.getElementsByTagName('script')[0]; s.parentNode.inser"))
	buf.Write(import_3.HTML("tBefore(ga, s);\n  })();\n</script>\n"))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}

func Header(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   1
	   column: 1
	*/
	buf.Write(import_3.HTML("\n  <header>\n\n  <div data-role=\"row\" class=\"row top\">\n\n    <div class=\"wrapper\">\n\n      <nav data-rol"))
	buf.Write(import_3.HTML("e=\"column\" class=\"column metanav\">\n        <ul>\n\n<li data-nav-key=\"page-4\" class=\"parent\">\n  <span>A"))
	buf.Write(import_3.HTML("bout orakel</span>\n\n\n\n    <ul>\n\n<li data-nav-key=\"page-5\" class=\"parent\">\n  <a href=\"/information/ab"))
	buf.Write(import_3.HTML("out-orakel/team\">Team</a>\n\n\n\n\n</li>\n\n    </ul>\n\n\n</li>\n\n<li data-nav-key=\"page-7\" class=\"parent\">\n  "))
	buf.Write(import_3.HTML("<a href=\"/information/faq\">Faq</a>\n\n\n\n\n</li>\n\n<li data-nav-key=\"page-8\" class=\"\">\n  <a href=\"/inform"))
	buf.Write(import_3.HTML("ation/contact\">Contact</a>\n\n\n\n\n</li>\n\n<li data-nav-key=\"page-209\" class=\"\">\n  <a href=\"/information/"))
	buf.Write(import_3.HTML("payments\">Payments</a>\n\n\n\n\n</li>\n\n        </ul>\n      </nav>\n\n      <nav data-role=\"column\" class=\"c"))
	buf.Write(import_3.HTML("olumn languages\">\n        <ul>\n          <li class=\"current\"><a href=\"#\" class=\"eu\" id=\"current-coun"))
	buf.Write(import_3.HTML("try-flag\">Belgium</a></li>\n\n          <li><a href=\"/\" data-nav-key=\"locale-en\">English</a></li>\n\n\n  "))
	buf.Write(import_3.HTML("      </ul>\n      </nav>\n\n      <div data-role=\"column\" class=\"column countries\" id=\"orakel-countrie"))
	buf.Write(import_3.HTML("s-list\">\n\n        <a href=\"http://be.orakel.com\" class=\"country be\">Belgium</a>\n\n        <a href=\"ht"))
	buf.Write(import_3.HTML("tp://orakel.com\" class=\"country eu\">Europe</a>\n\n        <a href=\"http://fr.orakel.com\" class=\"countr"))
	buf.Write(import_3.HTML("y fr\">France</a>\n\n        <a href=\"http://de.orakel.com\" class=\"country de\">Germany</a>\n\n        <a "))
	buf.Write(import_3.HTML("href=\"http://gr.orakel.com\" class=\"country gr\">Greece</a>\n\n        <a href=\"http://hu.orakel.com\" cl"))
	buf.Write(import_3.HTML("ass=\"country hu\">Hungary</a>\n\n        <a href=\"http://it.orakel.com\" class=\"country it\">Italy</a>\n\n "))
	buf.Write(import_3.HTML("       <a href=\"http://nl.orakel.com\" class=\"country nl\">Netherlands</a>\n\n        <a href=\"http://no"))
	buf.Write(import_3.HTML(".orakel.com\" class=\"country no\">Norway</a>\n\n        <a href=\"http://pl.orakel.com\" class=\"country pl"))
	buf.Write(import_3.HTML("\">Poland</a>\n\n        <a href=\"http://pt.orakel.com\" class=\"country pt\">Portugal</a>\n\n        <a hre"))
	buf.Write(import_3.HTML("f=\"http://es.orakel.com\" class=\"country es\">Spain</a>\n\n        <a href=\"http://tr.orakel.com\" class="))
	buf.Write(import_3.HTML("\"country tr\">Turkey</a>\n\n        <a href=\"http://uk.orakel.com\" class=\"country uk\">United Kingdom</a"))
	buf.Write(import_3.HTML(">\n\n      </div>\n\n    </div>\n\n  </div>\n\n  <div data-role=\"row\" class=\"row middle\">\n\n    <div class=\"w"))
	buf.Write(import_3.HTML("rapper\">\n\n      <div data-role=\"column\" class=\"column left\">\n\n          <a href=\"/\">Orakel</a>\n\n    "))
	buf.Write(import_3.HTML("  </div>\n\n      <div class=\"column right\">\n\n        <form action=\"/search\" method=\"get\" id=\"search\">"))
	buf.Write(import_3.HTML("\n          <input id=\"query\" name=\"query\" type=\"text\" />\n          <input name=\"commit\" type=\"submit"))
	buf.Write(import_3.HTML("\" value=\"Search\" />\n        </form>\n\n      </div>\n\n    </div>\n\n  </div>\n\n  <div data-role=\"row\" clas"))
	buf.Write(import_3.HTML("s=\"row tag\">\n\n    <div class=\"wrapper\">\n\n      <div data-role=\"column\" class=\"column left\">\n        "))
	buf.Write(import_3.HTML("Your Producer and Supplier for party and eventing material\n      </div>\n      <div data-role=\"column"))
	buf.Write(import_3.HTML("\" class=\"column right\">\n        Here to help, contact us at +32 (0) 14 38 80 80\n      </div>\n\n    </"))
	buf.Write(import_3.HTML("div>\n\n  </div>\n\n  <div data-role=\"row\" class=\"row bottom\">\n\n    <div class=\"wrapper\">\n      <nav>\n\n "))
	buf.Write(import_3.HTML("       <a href=\"/products/wristbands\" data-nav-key=\"page-13\">Wristbands</a>\n\n        <a href=\"/produ"))
	buf.Write(import_3.HTML("cts/tokens\" data-nav-key=\"page-42\">Tokens</a>\n\n        <a href=\"/products/number-bibs\" data-nav-key="))
	buf.Write(import_3.HTML("\"page-210\">Number Bibs</a>\n\n        <a href=\"/products/glow-sticks\" data-nav-key=\"page-147\">Glow Sti"))
	buf.Write(import_3.HTML("cks</a>\n\n        <a href=\"/products/coupons\" data-nav-key=\"page-74\">Coupons</a>\n\n        <a href=\"/p"))
	buf.Write(import_3.HTML("roducts/lanyards\" data-nav-key=\"page-79\">Lanyards</a>\n\n        <a href=\"/products/badges\" data-nav-k"))
	buf.Write(import_3.HTML("ey=\"page-71\">Badges</a>\n\n        <a href=\"/products/event-control-consumables\" data-nav-key=\"page-19"))
	buf.Write(import_3.HTML("1\">Event Control/Consumables</a>\n\n        <a href=\"/products/tickets\" data-nav-key=\"page-86\">Tickets"))
	buf.Write(import_3.HTML("</a>\n\n        <a href=\"/products/labels-and-displays\" data-nav-key=\"page-106\">Labels and Displays</a"))
	buf.Write(import_3.HTML(">\n\n        <a href=\"/products/ear-plugs\" data-nav-key=\"page-158\">Ear Plugs</a>\n\n        <a href=\"/pr"))
	buf.Write(import_3.HTML("oducts/balloons\" data-nav-key=\"page-152\">Balloons</a>\n\n        <a href=\"/products/cups\" data-nav-key"))
	buf.Write(import_3.HTML("=\"page-103\">Cups</a>\n\n        <a href=\"/products/playpool-balls\" data-nav-key=\"page-138\">Playpool Ba"))
	buf.Write(import_3.HTML("lls</a>\n\n      </nav>\n    </div>\n\n  </div>\n\n</header>\n"))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}

func Index(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   1
	   column: 1
	*/
	/*
	   Unhandled:
	     {{include}}
	*/
	/*
	   file:
	   line:   34
	   column: 12
	*/
	buf.Write(import_3.HTML("\n"))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}

func index_1(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   1
	   column: 40
	*/
	buf.Write(import_3.HTML("\n\n<div class=\"content\">\n  <ul>\n\n    "))
	/*
	   file:
	   line:   6
	   column: 5
	*/
	/*
	   Unhandled:
	     {{include}}
	*/
	/*
	   file:
	   line:   17
	   column: 14
	*/
	buf.Write(import_3.HTML("\n\n  </ul>\n\n  <h2>"))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_1 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_2 := Get(value_1, "title")
	/*
	   file:
	   line:   21
	   column: 7
	*/
	buf.Write(value_2)
	/*
	   file:
	   line:   21
	   column: 18
	*/
	buf.Write(import_3.HTML("</h2>\n  "))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_3 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_4 := Get(value_3, "body")
	/*
	   file:
	   line:   22
	   column: 3
	*/
	buf.Write(value_4)
	/*
	   file:
	   line:   22
	   column: 15
	*/
	buf.Write(import_3.HTML("\n\n</div>\n\n<aside>\n\n  "))
	/*
	   file:
	   line:   28
	   column: 3
	*/
	/*
	   Unhandled:
	     {{include}}
	*/
	/*
	   file:
	   line:   30
	   column: 12
	*/
	buf.Write(import_3.HTML("\n\n</aside>\n\n"))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}

func index_1_1(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   6
	   column: 24
	*/
	buf.Write(import_3.HTML("\n      <li>\n        <a href="))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_1 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_2 := import_1.ProductURL(value_1)
	/*
	   file:
	   line:   8
	   column: 17
	*/
	buf.Write(value_2)
	/*
	   file:
	   line:   8
	   column: 35
	*/
	buf.Write(import_3.HTML(" title=\"\">\n          "))
	/*
	   file:
	   line:   9
	   column: 11
	*/
	/*
	   Unhandled:
	     {{include}}
	*/
	/*
	   file:
	   line:   13
	   column: 18
	*/
	buf.Write(import_3.HTML("\n          <span>"))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_3 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_4 := Get(value_3, "name")
	/*
	   file:
	   line:   14
	   column: 17
	*/
	buf.Write(value_4)
	/*
	   file:
	   line:   14
	   column: 27
	*/
	buf.Write(import_3.HTML("</span>\n        </a>\n      </li>\n    "))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}

func index_1_1_1(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   9
	   column: 31
	*/
	buf.Write(import_3.HTML("\n            <img alt="))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_1 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_2 := Get(value_1, "name")
	/*
	   file:
	   line:   10
	   column: 22
	*/
	buf.Write(value_2)
	/*
	   file:
	   line:   10
	   column: 32
	*/
	buf.Write(import_3.HTML(" src="))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_3 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_4 := Get(value_3, "thumbnail")
	/*
	   file:
	   line:   10
	   column: 50
	*/
	value_5 := Get(value_4, "url")
	/*
	   file:
	   line:   10
	   column: 37
	*/
	buf.Write(value_5)
	/*
	   file:
	   line:   10
	   column: 56
	*/
	buf.Write(import_3.HTML(" />\n          "))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}

func index_1_1_2(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   11
	   column: 19
	*/
	buf.Write(import_3.HTML("\n            <img alt="))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_1 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_2 := Get(value_1, "name")
	/*
	   file:
	   line:   12
	   column: 22
	*/
	buf.Write(value_2)
	/*
	   file:
	   line:   12
	   column: 32
	*/
	buf.Write(import_3.HTML(" src=\"/missing.png\" />\n          "))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}

func index_1_2(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   28
	   column: 30
	*/
	buf.Write(import_3.HTML("\n    "))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_1 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_2 := Get(value_1, "template_name")
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_3 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_4 := render(value_2, value_3)
	/*
	   file:
	   line:   29
	   column: 5
	*/
	buf.Write(value_4)
	/*
	   file:
	   line:   29
	   column: 34
	*/
	buf.Write(import_3.HTML("\n  "))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}

func Layout(ctx import_2.Context, val import_2.Value) *import_3.Buffer {
	/*
	   file:
	   line:   0
	   column: 0
	*/
	buf := new(import_3.Buffer)
	/*
	   file:
	   line:   1
	   column: 1
	*/
	buf.Write(import_3.HTML("<!doctype html>\n<!--[if lt IE 7]> <html class=\"no-js lt-ie9 lt-ie8 lt-ie7\" lang=\"en\"> <![endif]-->\n<"))
	buf.Write(import_3.HTML("!--[if IE 7]>    <html class=\"no-js lt-ie9 lt-ie8\" lang=\"en\"> <![endif]-->\n<!--[if IE 8]>    <html c"))
	buf.Write(import_3.HTML("lass=\"no-js lt-ie9\" lang=\"en\"> <![endif]-->\n<!-- Consider adding a manifest.appcache: h5bp.com/d/Off"))
	buf.Write(import_3.HTML("line -->\n<!--[if gt IE 8]><!--> <html class=\"no-js\" lang=\"en\"> <!--<![endif]-->\n<head>\n  <meta chars"))
	buf.Write(import_3.HTML("et=\"utf-8\" />\n  <meta http-equiv=\"X-UA-Compatible\" content=\"IE=EmulateIE7,chrome=1\" />\n\n  <title>"))
	/*
	   file:
	   line:   11
	   column: 19
	*/
	var value_1 string = "title"
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_2 := yield(value_1)
	/*
	   file:
	   line:   11
	   column: 10
	*/
	buf.Write(value_2)
	/*
	   file:
	   line:   11
	   column: 29
	*/
	buf.Write(import_3.HTML("</title>\n\n\n  <meta name=\"viewport\" content=\"width=device-width,initial-scale=.8\" />\n  <link rel=\"sho"))
	buf.Write(import_3.HTML("rtcut icon\" type=\"image/x-icon\" href=\"/favicon.png\" />\n\n  <link rel=\"stylesheet\" href=\"/stylesheets/"))
	buf.Write(import_3.HTML("screen.css\" media=\"screen, projection\" />\n  <link rel=\"stylesheet\" href=\"/stylesheets/print.css\" med"))
	buf.Write(import_3.HTML("ia=\"print\" />\n  <!--[if IE]><link rel=\"stylesheet\" href=\"/stylesheets/ie.css\" media=\"screen, project"))
	buf.Write(import_3.HTML("ion\"><![endif]-->\n\n  <script src=\"/javascripts/modernizr-2.0.6.min.js\" charset=\"utf-8\"></script>\n\n  "))
	buf.Write(import_3.HTML("<script type=\"text/javascript\" src=\"http://use.typekit.com/ghv2mbk.js\"></script>\n  <script type=\"tex"))
	buf.Write(import_3.HTML("t/javascript\">try{Typekit.load();}catch(e){}</script>\n\n</head>\n<body data-active-nav-keys=\"locale-en"))
	buf.Write(import_3.HTML(" page-1\" class=\"default home home\">\n\n  "))
	/*
	   file:
	   line:   29
	   column: 12
	*/
	var value_3 string = "header"
	/*
	   file:
	   line:   29
	   column: 34
	*/
	var value_4 string = "title"
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_5 := yield(value_4)
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_6 := map[string]interface{}{"title": value_5}
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_7 := render(value_3, value_6)
	/*
	   file:
	   line:   29
	   column: 3
	*/
	buf.Write(value_7)
	/*
	   file:
	   line:   29
	   column: 45
	*/
	buf.Write(import_3.HTML("\n  <div id=\"wrapper\">\n    "))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	var value_8 = val
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_9 := Get(value_8, "yield")
	/*
	   file:
	   line:   31
	   column: 5
	*/
	buf.Write(value_9)
	/*
	   file:
	   line:   31
	   column: 16
	*/
	buf.Write(import_3.HTML("\n  </div>\n  "))
	/*
	   file:
	   line:   33
	   column: 13
	*/
	var value_10 string = "footer"
	/*
	   file:
	   line:   0
	   column: 0
	*/
	value_11 := render(value_10)
	/*
	   file:
	   line:   33
	   column: 3
	*/
	buf.Write(value_11)
	/*
	   file:
	   line:   33
	   column: 24
	*/
	buf.Write(import_3.HTML("\n\n</body>\n</html>\n"))
	/*
	   file:
	   line:   0
	   column: 0
	*/
	return buf
}
