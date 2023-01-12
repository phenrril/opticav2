<?php include_once "includes/header.php";
include "../conexion.php";
$id_user = $_SESSION['idUser'];
$permiso = "estadisticas";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
}




$usuarios = mysqli_query($conexion, "SELECT * FROM usuario");
$totalU= mysqli_num_rows($usuarios);
$clientes = mysqli_query($conexion, "SELECT * FROM cliente");
$totalC = mysqli_num_rows($clientes);
$productos = mysqli_query($conexion, "SELECT * FROM producto");
$totalP = mysqli_num_rows($productos);
$ventas = mysqli_query($conexion, "SELECT * FROM ventas");
$totalV = mysqli_num_rows($ventas);
?>
    <div class="d-sm-flex align-items-center justify-content-between mb-4">
        <h1 class="h3 mb-0 text-gray">Panel de Administración</h1>
    </div>

    <!-- Content Row -->
    <div class="row">
        <a class="col-xl-3 col-md-6 mb-4" href="usuarios.php">
            <div class="card border-left-primary shadow h-100 py-2 bg-warning">
                <div class="card-body">
                    <div class="row no-gutters align-items-center">
                        <div class="col mr-2">
                            <div class="text-xs font-weight-bold text-white text-uppercase mb-1">Usuarios</div>
                            <div class="h5 mb-0 font-weight-bold text-white"><?php echo $totalU; ?></div>
                        </div>
                        <div class="col-auto">
                            <i class="fas fa-user fa-2x text-gray-300"></i>
                        </div>
                    </div>
                </div>
            </div>
        </a>

        <!-- Earnings (Monthly) Card Example -->
        <a class="col-xl-3 col-md-6 mb-4" href="clientes.php">
            <div class="card border-left-success shadow h-100 py-2 bg-success">
                <div class="card-body">
                    <div class="row no-gutters align-items-center">
                        <div class="col mr-2">
                            <div class="text-xs font-weight-bold text-white text-uppercase mb-1">Clientes</div>
                            <div class="h5 mb-0 font-weight-bold text-white"><?php echo $totalC; ?></div>
                        </div>
                        <div class="col-auto">
                            <i class="fas fa-users fa-2x text-gray-300"></i>
                        </div>
                    </div>
                </div>
            </div>
        </a>

        <!-- Earnings (Monthly) Card Example -->
        <a class="col-xl-3 col-md-6 mb-4" href="productos.php">
            <div class="card border-left-info shadow h-100 py-2 bg-primary">
                <div class="card-body">
                    <div class="row no-gutters align-items-center">
                        <div class="col mr-2">
                            <div class="text-xs font-weight-bold text-white text-uppercase mb-1">Productos</div>
                            <div class="row no-gutters align-items-center">
                                <div class="col-auto">
                                    <div class="h5 mb-0 mr-3 font-weight-bold text-white"><?php echo $totalP; ?></div>
                                </div>
                                <div class="col">
                                    <div class="progress progress-sm mr-2">
                                        <div class="progress-bar bg-danger" role="progressbar" style="width: 50%" aria-valuenow="50" aria-valuemin="0" aria-valuemax="100"></div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="col-auto">
                            <i class="fas fa-clipboard-list fa-2x text-gray-300"></i>
                        </div>
                    </div>
                </div>
            </div>
        </a>

        <!-- Pending Requests Card Example -->
        <a class="col-xl-3 col-md-6 mb-4" href="ventas.php">
            <div class="card border-left-warning bg-danger shadow h-100 py-2">
                <div class="card-body">
                    <div class="row no-gutters align-items-center">
                        <div class="col mr-2">
                            <div class="text-xs font-weight-bold text-white text-uppercase mb-1">Ventas</div>
                            <div class="h5 mb-0 font-weight-bold text-white"><?php echo $totalV; ?></div>
                        </div>
                        <div class="col-auto">
                            <i class="fas fa-dollar-sign fa-2x text-white-300"></i>
                        </div>
                    </div>
                </div>
            </div>
        </a>
        <div class="col-lg-6">
            <div class="au-card m-b-30">
                <div class="au-card-inner">
                    <h3 class="title-2 m-b-40">Productos con stock mínimo</h3>
                    <canvas id="sales-chart"></canvas>
                </div>
            </div>
        </div>
        <div class="col-lg-6">
            <div class="au-card m-b-30">
                <div class="au-card-inner">
                    <h3 class="title-2 m-b-40">Productos más vendidos</h3>
                    <canvas id="polarChart"></canvas>
                </div>
            </div>
        </div>
        <div class="col-lg-6">
            <div class="au-card m-b-30">
                <div class="au-card-inner">
                    <h3 class="title-2 m-b-40">Calendario de Ventas</h3>
                    <form method="POST" id="111">
                    <input type="text" name="calendario" id="calendario">
                    <input type="button" id="fecha" name="fecha" value="mostrar fecha" onclick="search()" >
                    <div id="okcalendario">

<form>
  <label for="start-year">Start Year:</label><br>
  <select id="start-year" name="start-year">
  </select>
  <br>
  <label for="start-month">Start Month:</label><br>
  <select id="start-month" name="start-month">
    <option value="01">January</option>
    <option value="02">February</option>
    <option value="03">March</option>
    <option value="04">April</option>
    <option value="05">May</option>
    <option value="06">June</option>
    <option value="07">July</option>
    <option value="08">August</option>
    <option value="09">September</option>
    <option value="10">October</option>
    <option value="11">November</option>
    <option value="12">December</option>
  </select>
  <br>
  <label for="start-day">Start Day:</label><br>
  <select id="start-day" name="start-day">
  </select>
</form> 

<script>
  // Obtener los elementos select
  var yearSelect = document.getElementById("start-year");
  var daySelect = document.getElementById("start-day");

  // Fecha actual
  var currentDate = new Date();

  // Año actual
  var currentYear = currentDate.getFullYear();

  // Generar opciones para los próximos 10 años
  for (var i = 0; i < 10; i++) {
    // Año a mostrar en la opción
    var optionYear = currentYear + i;

    // Crear elemento option
    var option = document.createElement("option");

    // Establecer el valor y el texto de la opción
    option.value = optionYear;
    option.text = optionYear;

    // Añadir la opción al select
    yearSelect.add(option);
  }

  // Generar opciones para los días del mes actual
  var currentMonth = currentDate.getMonth();
  var daysInMonth = new Date(currentYear, currentMonth + 1, 0).getDate();
  for (var i = 1; i <= daysInMonth; i++) {
    // Crear elemento option
    var option = document.createElement("option");

    // Establecer el valor y el texto de la opción
    option.value = i;
    option.text = i;

    // Añadir la opción al select
    daySelect.add(option);
  };
</script>
<script>
   


  function search() {
    // Obtener los valores de los select
    var year = document.getElementById("start-year").value;
    var month = document.getElementById("start-month").value;
    var day = document.getElementById("start-day").value;

    // Generar la fecha en formato YYYY-MM-DD
    var startDate = year + "-" + month + "-" + day;

    // Tu consulta SQL
    var sql = "SELECT * FROM ventas WHERE date(fecha) = '" + startDate + "'";

    // Ejecutar la consulta
    executeQuery(sql);
    
  }
</script>


                    </div>
                    </form>
                    <canvas id="polarChart"></canvas>
                </div>
            </div>
        </div>
    </div>

<!-- <script>
    $('#fecha').click( function() {
    {$.ajax({
        url: "resultado2.php",
        type: "POST",
        data: $("#111").serialize(),
        success: function(resultado){
                $("#okcalendario").html(resultado);

                }
            });
        }
    });
</script> -->
<?php include_once "includes/footer.php"; ?>