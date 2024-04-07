package util

import (
	"fmt"
	"testing"
)

func TestGPXParser(t *testing.T) {
	data := `<?xml version="1.0" encoding="UTF-8"?>
	<gpx
	  version="1.0"
	  creator="GPSBabel - http://www.gpsbabel.org"
	  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	  xmlns="http://www.topografix.com/GPX/1/0"
	  xsi:schemaLocation="http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd">
	<time>2010-12-14T06:17:04Z</time>
	<bounds minlat="46.430350000" minlon="13.738842000" maxlat="46.435641000" maxlon="13.748333000"/>
	<trk>
	<trkseg>
	<trkpt lat="46.434981000" lon="13.748273000">
	  <ele>1614.678000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434890000" lon="13.748193000">
	  <ele>1636.776000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434911000" lon="13.748063000">
	  <ele>1632.935520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434871000" lon="13.747903000">
	  <ele>1631.502960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434951000" lon="13.747993000">
	  <ele>1629.095040</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435041000" lon="13.748043000">
	  <ele>1628.119680</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435511000" lon="13.748163000">
	  <ele>1626.199440</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435461000" lon="13.748003000">
	  <ele>1627.174800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434991000" lon="13.746923000">
	  <ele>1652.168400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434691000" lon="13.746203000">
	  <ele>1684.842960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434691000" lon="13.745963000">
	  <ele>1688.226240</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434721000" lon="13.745803000">
	  <ele>1688.226240</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434731000" lon="13.745563000">
	  <ele>1694.931840</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434761000" lon="13.745343000">
	  <ele>1699.747680</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434751000" lon="13.745183000">
	  <ele>1701.180240</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434721000" lon="13.745063000">
	  <ele>1702.155600</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434601000" lon="13.745053000">
	  <ele>1716.084960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434701000" lon="13.744773000">
	  <ele>1714.652400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434691000" lon="13.744593000">
	  <ele>1718.980560</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434321000" lon="13.743763000">
	  <ele>1725.716640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434301000" lon="13.743643000">
	  <ele>1726.173840</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434181000" lon="13.743282000">
	  <ele>1727.149200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433981000" lon="13.743142000">
	  <ele>1728.094080</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433271000" lon="13.741952000">
	  <ele>1738.670640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431170000" lon="13.740652000">
	  <ele>1874.215200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431060000" lon="13.740722000">
	  <ele>1883.359200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430920000" lon="13.740712000">
	  <ele>1887.687360</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430350000" lon="13.740702000">
	  <ele>1886.254800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430420000" lon="13.740582000">
	  <ele>1921.337280</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430590000" lon="13.740572000">
	  <ele>1959.315360</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430540000" lon="13.740452000">
	  <ele>1945.843200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430580000" lon="13.740352000">
	  <ele>1946.330880</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430630000" lon="13.740242000">
	  <ele>1951.603920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430700000" lon="13.740252000">
	  <ele>1954.987200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430740000" lon="13.740352000">
	  <ele>1956.419760</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430800000" lon="13.740262000">
	  <ele>1961.235600</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430860000" lon="13.740262000">
	  <ele>1968.428880</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430900000" lon="13.740142000">
	  <ele>1972.757040</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431000000" lon="13.740042000">
	  <ele>1973.244720</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431070000" lon="13.739972000">
	  <ele>1975.652640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431100000" lon="13.739872000">
	  <ele>1981.413360</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431280000" lon="13.739812000">
	  <ele>1992.965280</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431390000" lon="13.739732000">
	  <ele>1999.670880</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431460000" lon="13.739672000">
	  <ele>2000.646240</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431540000" lon="13.739752000">
	  <ele>2015.063280</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431580000" lon="13.739742000">
	  <ele>2019.391440</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431610000" lon="13.739682000">
	  <ele>2014.575600</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431640000" lon="13.739552000">
	  <ele>2019.879120</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431700000" lon="13.739452000">
	  <ele>2024.207280</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431720000" lon="13.739332000">
	  <ele>2025.152160</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431820000" lon="13.739242000">
	  <ele>2023.719600</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431930000" lon="13.739102000">
	  <ele>2028.992640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432000000" lon="13.739032000">
	  <ele>2032.375920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432221000" lon="13.739012000">
	  <ele>2057.369520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432221000" lon="13.738922000">
	  <ele>2050.145760</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432271000" lon="13.738842000">
	  <ele>2048.713200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432191000" lon="13.738842000">
	  <ele>2046.792960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432101000" lon="13.738902000">
	  <ele>2042.464800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432031000" lon="13.738972000">
	  <ele>2038.624320</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431970000" lon="13.739062000">
	  <ele>2034.753360</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431880000" lon="13.739112000">
	  <ele>2033.320800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431810000" lon="13.739182000">
	  <ele>2029.968000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431750000" lon="13.739282000">
	  <ele>2027.072400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431700000" lon="13.739382000">
	  <ele>2022.744240</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431660000" lon="13.739502000">
	  <ele>2021.799360</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431620000" lon="13.739602000">
	  <ele>2018.416080</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431610000" lon="13.739732000">
	  <ele>2016.983520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431520000" lon="13.739752000">
	  <ele>2015.550960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431440000" lon="13.739782000">
	  <ele>2008.814880</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431400000" lon="13.739722000">
	  <ele>2003.054160</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431290000" lon="13.739782000">
	  <ele>1997.750640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431200000" lon="13.739802000">
	  <ele>1999.213680</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431160000" lon="13.739822000">
	  <ele>1991.014560</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431110000" lon="13.739832000">
	  <ele>1983.333600</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431050000" lon="13.739922000">
	  <ele>1981.901040</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430950000" lon="13.739962000">
	  <ele>1982.845920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430930000" lon="13.740042000">
	  <ele>1979.005440</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430890000" lon="13.740182000">
	  <ele>1976.140320</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430860000" lon="13.740302000">
	  <ele>1972.757040</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430820000" lon="13.740262000">
	  <ele>1970.836800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430760000" lon="13.740172000">
	  <ele>1963.643520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430710000" lon="13.740252000">
	  <ele>1959.772560</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430670000" lon="13.740192000">
	  <ele>1956.907440</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430620000" lon="13.740272000">
	  <ele>1955.444400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430600000" lon="13.740402000">
	  <ele>1954.011840</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430560000" lon="13.740452000">
	  <ele>1948.738800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430510000" lon="13.740532000">
	  <ele>1946.818560</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.430530000" lon="13.740622000">
	  <ele>1943.922960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431020000" lon="13.741072000">
	  <ele>1942.490400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431130000" lon="13.741272000">
	  <ele>1942.490400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431260000" lon="13.741292000">
	  <ele>1942.490400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431310000" lon="13.741382000">
	  <ele>1942.490400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431430000" lon="13.741512000">
	  <ele>1942.490400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431490000" lon="13.741582000">
	  <ele>1942.490400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431560000" lon="13.741702000">
	  <ele>1942.490400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431630000" lon="13.741362000">
	  <ele>1873.270320</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431650000" lon="13.741252000">
	  <ele>1867.021920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431740000" lon="13.741282000">
	  <ele>1869.429840</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431800000" lon="13.741342000">
	  <ele>1863.181440</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431850000" lon="13.741442000">
	  <ele>1864.614000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.431920000" lon="13.741512000">
	  <ele>1864.614000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432000000" lon="13.741592000">
	  <ele>1864.614000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432060000" lon="13.741642000">
	  <ele>1864.614000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432140000" lon="13.741722000">
	  <ele>1864.614000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432220000" lon="13.741672000">
	  <ele>1852.604880</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432290000" lon="13.741622000">
	  <ele>1846.356480</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432330000" lon="13.741572000">
	  <ele>1837.700160</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432410000" lon="13.741622000">
	  <ele>1833.859680</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432500000" lon="13.741662000">
	  <ele>1829.531520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432580000" lon="13.741692000">
	  <ele>1827.123600</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432660000" lon="13.741702000">
	  <ele>1822.307760</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432750000" lon="13.741762000">
	  <ele>1819.899840</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432830000" lon="13.741822000">
	  <ele>1816.547040</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432870000" lon="13.741902000">
	  <ele>1813.651440</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432920000" lon="13.742012000">
	  <ele>1813.651440</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.432980000" lon="13.742092000">
	  <ele>1811.274000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433041000" lon="13.742192000">
	  <ele>1809.810960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433091000" lon="13.742272000">
	  <ele>1807.890720</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433141000" lon="13.742362000">
	  <ele>1805.482800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433191000" lon="13.742402000">
	  <ele>1800.697440</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433241000" lon="13.742462000">
	  <ele>1794.906240</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433301000" lon="13.742532000">
	  <ele>1792.041120</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433361000" lon="13.742612000">
	  <ele>1788.200640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433431000" lon="13.742692000">
	  <ele>1786.737600</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433511000" lon="13.742722000">
	  <ele>1782.897120</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433601000" lon="13.742712000">
	  <ele>1779.056640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433681000" lon="13.742732000">
	  <ele>1779.544320</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433701000" lon="13.742792000">
	  <ele>1772.320560</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433781000" lon="13.743022000">
	  <ele>1767.535200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433811000" lon="13.743132000">
	  <ele>1760.799120</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433841000" lon="13.743212000">
	  <ele>1759.823760</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433911000" lon="13.743292000">
	  <ele>1759.823760</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.433961000" lon="13.743392000">
	  <ele>1757.903520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434031000" lon="13.743472000">
	  <ele>1757.415840</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434121000" lon="13.743532000">
	  <ele>1755.495600</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434211000" lon="13.743542000">
	  <ele>1753.087680</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434271000" lon="13.743653000">
	  <ele>1751.655120</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434301000" lon="13.743753000">
	  <ele>1748.302320</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434311000" lon="13.743893000">
	  <ele>1747.326960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434321000" lon="13.743983000">
	  <ele>1742.053920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434341000" lon="13.744123000">
	  <ele>1737.238080</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434391000" lon="13.744233000">
	  <ele>1736.293200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434451000" lon="13.744313000">
	  <ele>1734.342480</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434541000" lon="13.744343000">
	  <ele>1732.909920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434581000" lon="13.744443000">
	  <ele>1730.989680</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434591000" lon="13.744523000">
	  <ele>1726.173840</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434631000" lon="13.744613000">
	  <ele>1723.765920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434641000" lon="13.744733000">
	  <ele>1719.468240</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434581000" lon="13.744853000">
	  <ele>1715.140080</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434591000" lon="13.744983000">
	  <ele>1714.164720</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434611000" lon="13.745113000">
	  <ele>1711.756800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434641000" lon="13.745233000">
	  <ele>1710.811920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434681000" lon="13.745363000">
	  <ele>1708.891680</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434731000" lon="13.745453000">
	  <ele>1708.404000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434711000" lon="13.745563000">
	  <ele>1704.563520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434681000" lon="13.745683000">
	  <ele>1702.643280</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434661000" lon="13.745783000">
	  <ele>1700.235360</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434631000" lon="13.745823000">
	  <ele>1696.852080</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434631000" lon="13.745933000">
	  <ele>1695.907200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434631000" lon="13.746043000">
	  <ele>1693.499280</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434591000" lon="13.746193000">
	  <ele>1692.523920</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434571000" lon="13.746293000">
	  <ele>1689.658800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434611000" lon="13.746383000">
	  <ele>1686.763200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434701000" lon="13.746473000">
	  <ele>1686.763200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434771000" lon="13.746593000">
	  <ele>1686.763200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434811000" lon="13.746733000">
	  <ele>1686.763200</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434921000" lon="13.746983000">
	  <ele>1686.275520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434941000" lon="13.747123000">
	  <ele>1686.275520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.434991000" lon="13.747253000">
	  <ele>1686.275520</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435151000" lon="13.747433000">
	  <ele>1685.818320</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435241000" lon="13.747543000">
	  <ele>1685.818320</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435291000" lon="13.747663000">
	  <ele>1685.330640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435381000" lon="13.747703000">
	  <ele>1685.330640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435471000" lon="13.747733000">
	  <ele>1685.330640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435531000" lon="13.747863000">
	  <ele>1685.330640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435601000" lon="13.747953000">
	  <ele>1684.842960</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435521000" lon="13.747783000">
	  <ele>1677.162000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435591000" lon="13.748093000">
	  <ele>1677.162000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435641000" lon="13.748223000">
	  <ele>1677.162000</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435501000" lon="13.748203000">
	  <ele>1659.849360</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435421000" lon="13.748233000">
	  <ele>1652.168400</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435431000" lon="13.748313000">
	  <ele>1649.272800</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435381000" lon="13.748333000">
	  <ele>1644.944640</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	<trkpt lat="46.435231000" lon="13.748253000">
	  <ele>1643.512080</ele>
	  <time>1901-12-13T20:45:52.2073437Z</time>
	</trkpt>
	</trkseg>
	</trk>
	</gpx>`

	gpxHandler, err := GPXParser([]byte(data))
	if err != nil {
		t.Log(err)
	}

	for _, track := range gpxHandler.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				fmt.Println(point.Latitude, point.Longitude)
			}
		}
	}
}